package sqs

import (
	"fmt"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	//"github.com/go-redsync/redsync/v4"
	//"github.com/go-redsync/redsync/v4/redis/redigo"
	//redredis "github.com/go-redsync/redsync/v4/redis
	//redigolib "github.com/gomodule/redigo/redis"

	"github.com/HomesNZ/go-common/redis"
	redsync "github.com/go-redsync/redsync" //TODO: replace it to new version
	redigo "github.com/gomodule/redigo/redis"
)

var contextLogger = logrus.WithField("package", "sqs_consumer")

const (
	maxMessages = 10

	maxRetries = 5

	defaultWaitSeconds = 10

	defaultVisibilityTimeout = 3

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

// AccessKeyID - is aws access key id
// SecretAccessKey - is aws secret access key
// QueueName - is aws SQS queue name
// Region - is aws SQS region
type Config struct {
	QueueName       string
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

type Consumer struct {
	config       Config
	conn         *sqs.SQS
	queueUrl     string
	handler      interface{}
	handlers     map[string]SNSMessageHandler
	responseChan chan *sqs.ReceiveMessageOutput
	doneChan     chan bool
	started      bool

	waitForCompletion bool

	redsyncEnabled bool
	redsync        *redsync.Redsync
	redsyncOptions []redsync.Option

	batchSize int
}

// ConfigFromEnv returns back the config with options from environment variables
func ConfigFromEnv() Config {
	return Config{
		AccessKeyID:     env.MustGetString("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: env.MustGetString("AWS_SECRET_ACCESS_KEY"),
		Region:          env.MustGetString("AWS_SQS_REGION"),
		QueueName:       env.MustGetString("AWS_SQS_QUEUE"),
	}
}

// NewConsumer returns a pointer to a fresh Consumer instance.
func NewConsumer(config Config, handlers map[string]SNSMessageHandler) (*Consumer, error) {
	if config.AccessKeyID == "" {
		return nil, errors.New("empty aws access key id")
	}
	if config.SecretAccessKey == "" {
		return nil, errors.New("empty aws secret access key")
	}
	if config.Region == "" {
		return nil, errors.New("empty aws sqs region")
	}

	sess := session.New(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     config.AccessKeyID,
			SecretAccessKey: config.SecretAccessKey,
		}}),
		MaxRetries: aws.Int(maxRetries),
	})

	s := sqs.New(sess)
	resultURL, err := s.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(config.QueueName),
	})
	if err != nil {
		contextLogger.Error("Can't get the SQS queue")
		return nil, err
	}

	router := NewRouter()
	for event, h := range handlers {
		router.AddRoute(event, h)
	}

	handler := Handler{
		Router: router,
	}

	return &Consumer{
		conn:      s,
		config:    config,
		queueUrl:  aws.StringValue(resultURL.QueueUrl),
		handler:   SNSMessageHandler(handler.HandleMessage),
		batchSize: maxMessages,
	}, nil
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

func (c *Consumer) BatchSize(size int) error {
	if c.started {
		return errors.New("BatchSize() called while consumer running")
	}

	c.batchSize = size
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

func (c *Consumer) Start() error {
	if c.started {
		return errors.New("can't start sqs consumer: already started")
	}

	c.responseChan = make(chan *sqs.ReceiveMessageOutput)
	c.doneChan = make(chan bool)
	c.started = true
	if c.redsyncEnabled {
		c.initRedsync()
	}

	go c.receive()
	go c.handleResponses()
	contextLogger.Info("now polling SQS queue:", c.config.QueueName)
	return nil
}

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
			c.started = false
			contextLogger.Info("stopped polling SQS queue:", c.config.QueueName)
			return
		default:
			contextLogger.Debug("waiting for request...")
			params := &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(c.queueUrl),
				MaxNumberOfMessages: aws.Int64(maxMessages),
				VisibilityTimeout:   aws.Int64(defaultVisibilityTimeout),
				WaitTimeSeconds:     aws.Int64(defaultWaitSeconds),
				MessageAttributeNames: aws.StringSlice([]string{
					"All",
				}),
			}
			resp, err := c.conn.ReceiveMessage(params)
			if err != nil {
				contextLogger.WithError(err).Errorf("Error occurred while receiving from SQS queue (%s), sleeping for %d seconds", err.Error(), secondsToSleepOnError)
				time.Sleep(time.Duration(secondsToSleepOnError) * time.Second)
				continue
			}
			contextLogger.Debug("request completed")
			c.responseChan <- resp
		}
	}
}

func (c Consumer) handleResponses() {
	for responce := range c.responseChan {
		wg := sync.WaitGroup{}
		for _, message := range responce.Messages {
			wg.Add(1)
			go func(message *sqs.Message) {
				defer wg.Done()
				c.handleMessage(*message)
			}(message)
			if c.waitForCompletion {
				wg.Wait()
			}
		}
		wg.Wait()
	}
}

func (c Consumer) handleMessage(message sqs.Message) {
	logger := contextLogger.WithFields(logrus.Fields{
		"receipt_handle": message.ReceiptHandle,
		"message_id":     *message.MessageId,
	})
	logger.Debug("handling message...")

	if c.handler == nil {
		// No handler supplied, don't handle!
		logger.Debug("no message handler supplied")
		return
	}

	//Lock this message in redsync
	if c.redsyncEnabled {
		name := redsyncPrefix + *message.MessageId
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
		if ok, err := handler(SQSMessage(message)); !ok {
			// Failed to handle message, do nothing. It's the responsibility of the
			// handler to communicate the failure via logs/bugsnag etc.
			logger.Debug("failed to handle message")
			logger.WithError(err).Error(err)
			return
		}
	case SNSMessageHandler:
		snsMessage, err := newSNSMessage(&message)
		if err != nil {
			logger.WithError(err).Error(err)
			return
		}
		var ok bool
		if ok, err = handler(snsMessage); !ok {
			// Failed to handle message, do nothing. It's the responsibility of the
			// handler to communicate the failure via logs/bugsnag etc.
			logger.Debug("failed to handle message")
			logger.WithError(err).Error(err)
			return
		}
	default:
		panic(fmt.Sprintf("Unknown handler: %v", c.handler))
	}

	logger.Debug("message handled, deleting...")
	// message was handled successfully, delete the message from SQS

	params := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueUrl),
		ReceiptHandle: message.ReceiptHandle,
	}
	_, err := c.conn.DeleteMessage(params)
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
	contextLogger.Info("stopping polling of SQS queue:", c.config.QueueName)
	c.doneChan <- true
	return nil
}
