package sqs_v2

import (
	"context"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/sqs_v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/sirupsen/logrus"
)

const (
	defaultWaitSeconds    = 10
	secondsToSleepOnError = 10
	maxRetries            = 5
)

type MessageHandler func(ctx context.Context, message []Message) error

type Consumer interface {
	Start(ctx context.Context)
	Stop()
}

type consumer struct {
	client   *SQS
	config   *config.Config
	log      *logrus.Entry
	handler  MessageHandler
	doneChan chan bool
	queueUrl *string
}

func (c *consumer) Start(ctx context.Context) {
	wg := &sync.WaitGroup{}
	wg.Add(c.config.MaxWorker)
	c.doneChan = make(chan bool)
	c.log.Info("now polling SQS queue:", c.config.QueueName)
	for i := 0; i < c.config.MaxWorker; i++ {
		go c.worker(ctx, wg)
	}
}

// Stop sends true to the doneChan, which stops the long polling process. Has to
// wait for the current poll to complete before the polling is stopped.
func (c consumer) Stop() {
	c.log.Info("stopping polling of SQS queue:", c.config.QueueName)
	c.doneChan <- true
}

func (c *consumer) worker(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-c.doneChan:
			close(c.doneChan)
			c.log.Info("stopped polling SQS queue:", c.config.QueueName)
			return
		case <-ctx.Done():
			c.log.Info("stopped polling SQS queue:", c.config.QueueName)
			return
		default:
			msgs, err := c.client.Receive(ctx, c.config.QueueName, defaultWaitSeconds, c.config.MaxMsg)
			if err != nil {
				c.log.WithError(err).Errorf("Error occurred while receiving from SQS queue (%s), sleeping for %d seconds", err.Error(), secondsToSleepOnError)
				time.Sleep(time.Duration(secondsToSleepOnError) * time.Second)
				continue
			}

			if len(msgs) == 0 {
				continue
			}
			c.consume(ctx, msgs)
			//c.async(ctx, msgs)
		}
	}
}

func (c consumer) async(ctx context.Context, msgs []types.Message) {
	wg := &sync.WaitGroup{}
	go func(m []types.Message) {
		defer wg.Done()
		c.consume(ctx, m)
	}(msgs)
	wg.Wait()
}

func (c consumer) consume(ctx context.Context, msgs []types.Message) {
	messages := make([]Message, 0, len(msgs))
	for _, m := range msgs {
		msg, err := newMessage(m)
		if err != nil {
			c.log.WithError(err).Error("failed to convert message")
		}
		messages = append(messages, msg)
	}
	if err := c.handler(ctx, messages); err != nil {
		// Failed to handle message, do nothing. It's the responsibility of the
		// handler to communicate the failure via logs/bugsnag etc.
		c.log.Debug("failed to handle message")
		c.log.WithError(err).Error(err)
	}

	for _, msg := range msgs {
		if err := c.client.Delete(ctx, c.config.QueueName, *msg.ReceiptHandle); err != nil {
			c.log.WithError(err).Error(err)
		}
	}
}