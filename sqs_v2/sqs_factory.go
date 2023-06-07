package sqs_v2

import (
	"context"
	"time"

	"github.com/HomesNZ/go-common/sqs_v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(msg string)
	Infof(msg string, args ...interface{})
	Error(err error, msg string)
}

type Options func(*Consumer)

// WithLogger sets the logger for the consumer
func WithLogger(logger Logger) Options {
	return func(c *Consumer) {
		c.log = logger
	}
}

func NewFromEnv(ctx context.Context, log *logrus.Entry, handler MessageHandler, options ...Options) (*Consumer, error) {
	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newConsumer(ctx, config, log, handler, options...)
}

// New returns a pointer to a fresh Consumer instance.
func newConsumer(ctx context.Context, config *config.Config, log *logrus.Entry, handler MessageHandler, options ...Options) (*Consumer, error) {
	cfg, err := awsCfg.LoadDefaultConfig(ctx, awsCfg.WithRegion(config.Region), awsCfg.WithRetryer(func() aws.Retryer {
		return retry.AddWithMaxAttempts(retry.NewStandard(), maxRetries)
	}))
	if err != nil {
		return nil, err
	}

	s := sqs.NewFromConfig(cfg)

	resultURL, err := s.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(config.QueueName),
	})
	if err != nil {
		return nil, err
	}

	sqsClient := &SQS{
		client:  s,
		timeout: time.Second * 5,
	}

	consumer := &Consumer{
		client:   sqsClient,
		config:   config,
		queueUrl: resultURL.QueueUrl,
		handler:  handler,
	}

	for _, opt := range options {
		opt(consumer)
	}

	return consumer, nil
}
