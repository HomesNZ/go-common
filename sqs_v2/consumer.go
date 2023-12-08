package sqs_v2

import (
	"context"
	"errors"
	"fmt"
	"github.com/HomesNZ/go-common/trace"
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

type MessageHandler func(ctx context.Context, message Message) error
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
		c.log.Info(ctx, "start polling", "queue_name", c.config.QueueName)
	}
	for i := 0; i < c.config.MaxWorker; i++ {
		go c.worker(ctx, wg)
	}
}

// Stop sends true to the doneChan, which stops the long polling process. Has to
// wait for the current poll to complete before the polling is stopped.
func (c *Consumer) Stop(ctx context.Context) {
	if c.log == nil {
		c.log.Info(ctx, "stop polling", "queue_name", c.config.QueueName)
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
				c.log.Info(ctx, "stopped polling", "queue_name", c.config.QueueName)
			}
			return
		//case <-ctx.Done():
		//	c.log.Info("stopped polling SQS queue:", c.config.QueueName)
		//	return
		default:
			msgs, err := c.client.Receive(ctx, c.config.QueueName, defaultWaitSeconds, c.config.MaxMsg)
			if err != nil {
				msg := fmt.Sprintf("error occurred while receiving from SQS queue (%s), sleeping for %d seconds", err.Error(), secondsToSleepOnError)
				if c.notifier != nil {
					c.notifier(errors.New(msg))
				}
				if c.log == nil {
					c.log.Error(ctx, msg)
				}
				time.Sleep(time.Duration(secondsToSleepOnError) * time.Second)
				continue
			}
			if c.log != nil {
				c.log.Info(ctx, "pulled messages", "len", len(msgs))
			}
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

func (c *Consumer) consume(ctx context.Context, messages []types.Message) {
	if ctx == nil {
		ctx = context.Background()
	}

	sem := make(chan struct{}, c.config.MaxHandlers)
	var wg sync.WaitGroup
	wg.Add(len(messages))
	for idx := range messages {
		sem <- struct{}{}

		go func(sqsMsg types.Message) {
			defer func() {
				<-sem
				wg.Done()
			}()

			timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(c.config.MaxMessageHandleTime))
			defer cancel()

			done := make(chan struct{})

			go func() {
				homesMessage, err := newMessage(sqsMsg)
				if err != nil && c.log != nil {
					c.log.Error(timeoutCtx, "failed to convert message", "reason", err.Error())
				}

				tracedCtx := trace.LinkCtxFromTrace(timeoutCtx, homesMessage.Trace)

				if err := c.handler(tracedCtx, homesMessage); err != nil && c.log != nil {
					// Failed to handle message, do nothing. It's the responsibility of the
					// handler to communicate the failure via logs/bugsnag etc.
					c.log.Error(tracedCtx, "failed to handle message", "reason", err.Error())
				}
				if err := c.client.Delete(tracedCtx, c.config.QueueName, *sqsMsg.ReceiptHandle); err != nil && c.log != nil {
					c.log.Error(tracedCtx, "failed to delete message", "reason", err.Error())
				}

				close(done)
			}()

			select {
			case <-done:
				return
			case <-timeoutCtx.Done():
				if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
					if c.log != nil {
						c.log.Error(timeoutCtx, "timeout exceeded", "reason", timeoutCtx.Err().Error())
					}
				}
			}
		}(messages[idx])
	}
	wg.Wait()
}
