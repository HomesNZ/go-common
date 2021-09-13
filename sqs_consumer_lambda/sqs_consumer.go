package sqsConsumerLambda

import (
	"fmt"
	"sync"
	"github.com/aws/aws-lambda-go/events"
	"github.com/sirupsen/logrus"
)

var contextLogger = logrus.WithField("package", "sqs_consumer_lambda")

type Consumer interface {
	Handle(sqsEvent events.SQSEvent)
}

type consumer struct {
	handler           interface{}
	waitForCompletion bool
}

func (c *consumer) Handle(sqsEvent events.SQSEvent) {
	wg := sync.WaitGroup{}
	for _, record := range sqsEvent.Records {
		wg.Add(1)
		go func(message *events.SQSMessage) {
			defer wg.Done()
			c.handleMessage(*message)
		}(&record)
		if c.waitForCompletion {
			wg.Wait()
		}
		wg.Wait()
	}
}

func (c consumer) handleMessage(message events.SQSMessage) {
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
	switch handler := c.handler.(type) {
	case SQSMessageHandler:
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
}