package sqsConsumerLambda

import (
	"fmt"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/redis"
	"github.com/aws/aws-lambda-go/events"
	"github.com/go-redsync/redsync"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// redsyncPrefix is the prefix added to the redsync key (to prevent multiple
	// processing of the same message).
	redsyncPrefix = "sqs:message:"

	// redsyncDefaultExpiry is the default duration redsync will lock a message
	// for. Can be overridden using Consumer.RedsyncOptions().
	redsyncDefaultExpiry = time.Second * 120
)

var contextLogger = logrus.WithField("package", "sqs_consumer_lambda")

type Consumer struct {
	handler           interface{}
	waitForCompletion bool
	redsyncEnabled    bool
	redsync           *redsync.Redsync
	redsyncOptions    []redsync.Option
}

// NewConsumer for AWS Lambda
func NewConsumer(rd redis.Cache, handlers map[string]SNSMessageHandler) (*Consumer, error) {
	if len(handlers) == 0 {
		return nil, errors.New("no handlers provided")
	}

	router := NewRouter()
	for event, h := range handlers {
		router.AddRoute(event, h)
	}

	handler := Handler{
		Router: router,
	}

	redsync := redsync.New(
		[]redsync.Pool{rd.Pool},
	)

	return &Consumer{
		handler:           SNSMessageHandler(handler.HandleMessage),
		waitForCompletion: true,
		redsyncEnabled:    true,
		redsync:           redsync,
	}, nil
}

func (c *Consumer) Handle(sqsEvent events.SQSEvent) {
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

func (c Consumer) handleMessage(message events.SQSMessage) {
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

	//Lock this message in redsync
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

func (c Consumer) redsyncDefaultOptions() []redsync.Option {
	return []redsync.Option{
		redsync.SetExpiry(redsyncDefaultExpiry),
		redsync.SetTries(1), // only try to lock once, then give up
	}
}
