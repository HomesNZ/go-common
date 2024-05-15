package sqs_v2

import (
	"context"
	"time"

	"github.com/HomesNZ/go-common/sqs_v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
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

func WithCredentials(key, secret, session string) Options {
	return func(c *Consumer) {
		c.config.AwsKey = key
		c.config.AwsSecret = secret
		c.config.AwsSession = session
	}
}

func NewFromEnv(ctx context.Context, handler MessageHandler, options ...Options) (*Consumer, error) {
	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newConsumer(ctx, config, handler, options...)
}

// New returns a pointer to a fresh Consumer instance.
func newConsumer(ctx context.Context, config *config.Config, handler MessageHandler, options ...Options) (*Consumer, error) {
	consumer := &Consumer{}
	consumer.config = config

	for _, opt := range options {
		opt(consumer)
	}
	var cfg aws.Config
	var err error
	if consumer.config.AwsSession != "" && consumer.config.AwsKey != "" && consumer.config.AwsSecret != "" {
		cfg, err = awsCfg.LoadDefaultConfig(ctx,
			awsCfg.WithRegion(config.Region),
			awsCfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(consumer.config.AwsKey, consumer.config.AwsSecret, consumer.config.AwsSession)),
			awsCfg.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), maxRetries)
			}))
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = awsCfg.LoadDefaultConfig(ctx,
			awsCfg.WithRegion(config.Region),
			awsCfg.WithRetryer(func() aws.Retryer {
				return retry.AddWithMaxAttempts(retry.NewStandard(), maxRetries)
			}))
		if err != nil {
			return nil, err
		}

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

	consumer.client = sqsClient
	consumer.queueUrl = resultURL.QueueUrl
	consumer.handler = handler

	return consumer, nil
}
