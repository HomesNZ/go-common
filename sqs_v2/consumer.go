package sqs_v2

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/sqs_v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	defaultWaitSeconds    = 10
	secondsToSleepOnError = 10
	maxRetries            = 5
)

type MessageHandler func(ctx context.Context, message []Message) error
type Notifier func(err error, rawData ...interface{})

type Consumer struct {
	client   *SQS
	config   *config.Config
	handler  MessageHandler
	doneChan chan bool
	queueUrl *string
	notifier Notifier
	log      Logger
}

func (c *Consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.config.MaxWorker)
	c.doneChan = make(chan bool)
	if c.log == nil {
		c.log.Infof("now polling SQS queue: %s", c.config.QueueName)
	}
	for i := 0; i < c.config.MaxWorker; i++ {
		go c.worker(ctx, wg)
	}
}

// Stop sends true to the doneChan, which stops the long polling process. Has to
// wait for the current poll to complete before the polling is stopped.
func (c *Consumer) Stop() {
	if c.log == nil {
		c.log.Infof("stopping polling of SQS queue: %s", c.config.QueueName)
	}
	c.doneChan <- true
}

func (c *Consumer) SetNotifier(f Notifier) {
	c.notifier = f
}

func (c *Consumer) worker(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-c.doneChan:
			close(c.doneChan)
			if c.log == nil {
				c.log.Infof("stopped polling SQS queue: %s", c.config.QueueName)
			}
			return
		//case <-ctx.Done():
		//	c.log.Info("stopped polling SQS queue:", c.config.QueueName)
		//	return
		default:
			msgs, err := c.client.Receive(ctx, c.config.QueueName, defaultWaitSeconds, c.config.MaxMsg)
			if err != nil {
				msg := fmt.Sprintf("Error occurred while receiving from SQS queue (%s), sleeping for %d seconds", err.Error(), secondsToSleepOnError)
				if c.notifier != nil {
					c.notifier(errors.New(msg))
				}
				if c.log == nil {
					c.log.Error(err, msg)
				}
				time.Sleep(time.Duration(secondsToSleepOnError) * time.Second)
				continue
			}
			c.log.Infof("pulled %d messages", len(msgs))
			if len(msgs) == 0 {
				continue
			}
			c.consume(ctx, msgs)
			//c.async(ctx, msgs)
		}
	}
}

func (c *Consumer) async(ctx context.Context, msgs []types.Message) {
	wg := &sync.WaitGroup{}
	go func(m []types.Message) {
		defer wg.Done()
		c.consume(ctx, m)
	}(msgs)
	wg.Wait()
}

func (c *Consumer) consume(ctx context.Context, msgs []types.Message) {
	messages := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		msg, err := newMessage(m)
		if err != nil && c.log != nil {
			c.log.Error(err, "failed to convert message")
		}
		messages = append(messages, msg)
	}
	if err := c.handler(ctx, messages); err != nil && c.log != nil {
		// Failed to handle message, do nothing. It's the responsibility of the
		// handler to communicate the failure via logs/bugsnag etc.
		c.log.Error(err, "failed to handle message")
		return
	}

	for _, msg := range msgs {
		if err := c.client.Delete(ctx, c.config.QueueName, *msg.ReceiptHandle); err != nil && c.log != nil {
			c.log.Error(err, "failed to delete message")
		}
	}
}
