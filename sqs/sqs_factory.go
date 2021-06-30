package sqs

import (
	"github.com/HomesNZ/go-common/sqs/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func NewFromEnv(handlers map[string]SNSMessageHandler) (Consumer, error) {
	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newConsumer(config, handlers)
}

func New(config config.Config, handlers map[string]SNSMessageHandler) (Consumer, error) {
	return newConsumer(config, handlers)
}

// New returns a pointer to a fresh Consumer instance.
func newConsumer(config config.Config, handlers map[string]SNSMessageHandler) (Consumer, error) {

	sess := session.New(&aws.Config{
		Region: aws.String(config.Region()),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     config.AccessKeyID(),
			SecretAccessKey: config.SecretAccessKey(),
		}}),
		MaxRetries: aws.Int(maxRetries),
	})

	s := sqs.New(sess)
	resultURL, err := s.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(config.QueueName()),
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
		queueUrl:  aws.StringValue(resultURL.QueueUrl),
		handler:   SNSMessageHandler(handler.HandleMessage),
		batchSize: maxMessages,
	}, nil
}
