package sqs

import (
	"context"
	"github.com/HomesNZ/go-common/sqs/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/sirupsen/logrus"
	"time"
)

func NewFromEnv(ctx context.Context, log *logrus.Entry, handler MessageHandler) (Consumer, error) {
	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newConsumer(ctx, config, log, handler)
}

// New returns a pointer to a fresh Consumer instance.
func newConsumer(ctx context.Context, config *config.Config, log *logrus.Entry, handler MessageHandler) (Consumer, error) {
	s := sqs.NewFromConfig(aws.Config{
		Region:           config.Region,
		Credentials:      credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, ""),
		RetryMaxAttempts: maxRetries,
	})

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

	return &consumer{
		client:   sqsClient,
		config:   config,
		queueUrl: resultURL.QueueUrl,
		handler:  handler,
		log:      log,
	}, nil
}
