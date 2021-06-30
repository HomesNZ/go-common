package sqs

import (
	"fmt"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/sqs/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var contextLogger = logrus.WithField("package", "sqs_consumer")

const (
	maxMessages              = 10
	maxRetries               = 5
	defaultWaitSeconds       = 10
	defaultVisibilityTimeout = 1800 // 30 mins -- prevent other consumers from processing the message again
	// secondsToSleepOnError defines the number of seconds to sleep for when an
	// error occurs while reciving SQS messages.
	secondsToSleepOnError = 10
)

type Consumer interface {
	Start() error
	BatchSize(size int) error
	WaitForCompletion(b bool)
	Stop() error
}

type consumer struct {
	config            config.Config
	conn              *sqs.SQS
	queueUrl          string
	handler           interface{}
	handlers          map[string]SNSMessageHandler
	responseChan      chan *sqs.ReceiveMessageOutput
	doneChan          chan bool
	started           bool
	waitForCompletion bool
	batchSize         int
}

func (c *consumer) BatchSize(size int) error {
	if c.started {
		return errors.New("BatchSize() called while consumer running")
	}

	c.batchSize = size
	return nil
}

// WaitForCompletion will make the consumer wait for each batch of messages to
// finish processing before it requests the next batch.
func (c *consumer) WaitForCompletion(b bool) {
	c.waitForCompletion = b
}

func (c *consumer) Start() error {
	if c.started {
		return errors.New("can't start sqs consumer: already started")
	}

	c.responseChan = make(chan *sqs.ReceiveMessageOutput)
	c.doneChan = make(chan bool)
	c.started = true

	go c.receive()
	go c.handleResponses()
	contextLogger.Info("now polling SQS queue:", c.config.QueueName())
	return nil
}

func (c consumer) receive() {
	for {
		select {
		case <-c.doneChan:
			close(c.doneChan)
			close(c.responseChan)
			c.doneChan = nil
			c.responseChan = nil
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

func (c consumer) handleResponses() {
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

func (c consumer) handleMessage(message sqs.Message) {
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
func (c consumer) Stop() error {
	if !c.started {
		return errors.New("can't stop sqs consumer: already stopped")
	}
	contextLogger.Info("stopping polling of SQS queue:", c.config.QueueName)
	c.doneChan <- true
	return nil
}
