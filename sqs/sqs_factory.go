package sqs

import (
	"context"
	"github.com/HomesNZ/go-common/sqs/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func NewFromEnv(ctx context.Context, handlers map[string]SNSMessageHandler) (Consumer, error) {
	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newConsumer(ctx, config, handlers)
}

func New(ctx context.Context, config *config.Config, handlers map[string]SNSMessageHandler) (Consumer, error) {
	return newConsumer(ctx, config, handlers)
}

// New returns a pointer to a fresh Consumer instance.
func newConsumer(ctx context.Context, config *config.Config, handlers map[string]SNSMessageHandler) (Consumer, error) {
	s := sqs.NewFromConfig(aws.Config{
		Region:           config.Region,
		Credentials:      credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, ""),
		RetryMaxAttempts: maxRetries,
	})

	resultURL, err := s.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
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

	return &consumer{
		conn:      s,
		config:    config,
		queueUrl:  *resultURL.QueueUrl,
		handler:   SNSMessageHandler(handler.HandleMessage),
		batchSize: maxMessages,
	}, nil
}
