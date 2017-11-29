package sqsConsumer

import (
	"errors"
	"fmt"
	"sync"
	"time"

	redsync "gopkg.in/redsync.v1"

	"github.com/HomesNZ/go-common/redis"
	"github.com/Sirupsen/logrus"
	"github.com/goamz/goamz/sqs"

	redigo "github.com/garyburd/redigo/redis"
)

// TODO: add tests

var contextLogger = logrus.WithField("package", "sqs_consumer")

const (
	maxMessages = 10

	// secondsToSleepOnError defines the number of seconds to sleep for when an
	// error occurs while reciving SQS messages.
	secondsToSleepOnError = 10

	// redsyncPrefix is the prefix added to the redsync key (to prevent multiple
	// processing of the same message).
	redsyncPrefix = "sqs:message:"

	// redsyncDefaultExpiry is the default duration redsync will lock a message
	// for. Can be overridden using Consumer.RedsyncOptions().
	redsyncDefaultExpiry = time.Second * 120
)

type ConsumerMessage interface {
	changeMessageVisibility(Consumer, time.Duration) error
}

// Consumer contains all the channels to manage goroutines and the SQS
// connection.
type Consumer struct {
	conn         *sqs.SQS
	queueName    string
	queue        *sqs.Queue
	handler      interface{}
	responseChan chan *sqs.ReceiveMessageResponse
	doneChan     chan bool
	started      bool

	waitForCompletion bool

	redsyncEnabled bool
	redsync        *redsync.Redsync
	redsyncOptions []redsync.Option
}

// NewConsumer returns a pointer to a fresh Consumer instance.
func NewConsumer(conn *sqs.SQS, queueName string, handler interface{}) *Consumer {
	return &Consumer{
		conn:      conn,
		queueName: queueName,
		handler:   handler,
	}
}

// RedsyncEnabled uses redsync to prevent multiple processing of the same SQS
// message.
func (c *Consumer) RedsyncEnabled(b bool) {
	if c.started {
		contextLogger.Error("RedsyncEnabled() called while consumer running")
		return
	}
	c.redsyncEnabled = b
}

// RedsyncOptions sets custom options for Redsync.
func (c *Consumer) RedsyncOptions(options []redsync.Option) {
	if c.started {
		contextLogger.Error("RedsyncOptions() called while consumer running")
		return
	}
	c.redsyncOptions = options
}

func (c Consumer) redsyncDefaultOptions() []redsync.Option {
	return []redsync.Option{
		redsync.SetExpiry(redsyncDefaultExpiry),
		redsync.SetTries(1), // only try to lock once, then give up
	}
}

// WaitForCompletion will make the consumer wait for each batch of messages to
// finish processing before it requests the next batch.
func (c *Consumer) WaitForCompletion(b bool) {
	c.waitForCompletion = b
}

// Start attempts to initialize the long polling process.
func (c *Consumer) Start() error {
	if c.started {
		return errors.New("can't start sqs consumer: already started")
	}
	q, err := c.conn.GetQueue(c.queueName)
	if err != nil {
		return err
	}
	c.queue = q
	c.responseChan = make(chan *sqs.ReceiveMessageResponse)
	c.doneChan = make(chan bool)
	c.started = true
	if c.redsyncEnabled {
		c.initRedsync()
	}
	go c.receive()
	go c.handleResponses()
	contextLogger.Info("now polling SQS queue:", c.queueName)
	return nil
}

// RedisPool is a redis pool wrapper for redsync
type RedisPool struct{}

// Get implements redsync.Pool
func (r RedisPool) Get() redigo.Conn {
	return redis.CacheConn().Conn()
}

func (c *Consumer) initRedsync() {
	p := RedisPool{}
	c.redsync = redsync.New(
		[]redsync.Pool{p},
	)
}

func (c *Consumer) terminateRedsync() {
	c.redsync = nil
}

// recieve handles the SQS long polling process. It passes messages as it
// recieves them to the responseChan for handleResponses to handle them. If it
// receives a message on the doneChan, it'll close all channels, log a message,
// and end the goroutine.
func (c Consumer) receive() {
	for {
		select {
		case <-c.doneChan:
			close(c.doneChan)
			close(c.responseChan)
			c.doneChan = nil
			c.responseChan = nil

			if c.redsyncEnabled {
				c.terminateRedsync()
			}

			c.queue = nil
			c.started = false
			contextLogger.Info("stopped polling SQS queue:", c.queueName)
			return
		default:
			contextLogger.Debug("waiting for request...")
			response, err := c.queue.ReceiveMessage(maxMessages)
			if err != nil {
				contextLogger.WithError(err).Errorf("Error occurred while receiving from SQS queue (%s), sleeping for %d seconds", err.Error(), secondsToSleepOnError)
				time.Sleep(time.Duration(secondsToSleepOnError) * time.Second)
				continue
			}
			contextLogger.Debug("request completed")
			c.responseChan <- response
		}
	}
}

// handleResponses blocks its goroutine waiting for records to come through on
// the response channel. As it receives responses, it spawns a goroutine to
// handle each message recieved (1:1). This means messages are not necessarily
// handled in order.
func (c Consumer) handleResponses() {
	for response := range c.responseChan {
		wg := sync.WaitGroup{}
		wg.Add(len(response.Messages))
		for _, message := range response.Messages {
			go func(message sqs.Message) {
				defer wg.Done()
				c.handleMessage(message)
			}(message)
		}
		if c.waitForCompletion {
			wg.Wait()
		}
	}
}

// handleMessage passes the handling off to c.handler. If the message is handled
// successfully by c.handler, then it is deleted from SQS.
func (c Consumer) handleMessage(message sqs.Message) {
	logger := contextLogger.WithFields(logrus.Fields{
		"receipt_handle": message.ReceiptHandle,
		"message_id":     message.MessageId,
	})
	logger.Debug("handling message...")

	if c.handler == nil {
		// No handler supplied, don't handle!
		logger.Debug("no message handler supplied")
		return
	}

	// Lock this message in redsync
	if c.redsyncEnabled {
		name := redsyncPrefix + message.MessageId
		options := c.redsyncDefaultOptions()
		if c.redsyncOptions != nil {
			options = append(options, c.redsyncOptions...)
		}

		mutex := c.redsync.NewMutex(name, options...)
		err := mutex.Lock()
		if err != nil {
			logger.Warn("can't acquire redsync lock, refusing to handle message (duplicate?): ", err)
			return
		}

		defer mutex.Unlock()
	}

	switch handler := c.handler.(type) {
	case MessageHandler:
		if !handler(SQSMessage(message)) {
			// Failed to handle message, do nothing. It's the responsibility of the
			// handler to communicate the failure via logs/bugsnag etc.
			logger.Debug("failed to handle message")
			return
		}
	case SNSMessageHandler:
		snsMessage, err := newSNSMessage(&message)
		if err != nil {
			logger.WithError(err).Error(err)
			return
		}

		if !handler(snsMessage) {
			// Failed to handle message, do nothing. It's the responsibility of the
			// handler to communicate the failure via logs/bugsnag etc.
			logger.Debug("failed to handle message")
			return
		}
	default:
		panic(fmt.Sprintf("Unknown handler: %v", c.handler))
	}

	logger.Debug("message handled, deleting...")
	// message was handled successfully, delete the message from SQS
	_, err := c.queue.DeleteMessage(&message)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug("message deleted")
}

// Stop sends true to the doneChan, which stops the long polling process. Has to
// wait for the current poll to complete before the polling is stopped.
func (c Consumer) Stop() error {
	if !c.started {
		return errors.New("can't stop sqs consumer: already stopped")
	}
	contextLogger.Info("stopping polling of SQS queue:", c.queueName)
	c.doneChan <- true
	return nil
}

// ChangeMessageVisibility sets the visibility timeout for a given message.
// http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html
func (c Consumer) ChangeMessageVisibility(message ConsumerMessage, d time.Duration) error {
	return message.changeMessageVisibility(c, d)
}
