package sqs_v2

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/pkg/errors"
	"time"
)

type SQS struct {
	client  *sqs.Client
	timeout time.Duration
}

func (s SQS) Receive(ctx context.Context, queueURL string, waitTimeSeconds int32, maxMsg int32) ([]types.Message, error) {
	if maxMsg < 1 || maxMsg > 10 {
		return nil, errors.New("msgMax valid values: 1 to 10")
	}
	res, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   maxMsg,
		WaitTimeSeconds:       waitTimeSeconds,
		MessageAttributeNames: []string{"All"},
		//VisibilityTimeout:
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	return res.Messages, nil
}

func (s SQS) Delete(ctx context.Context, queueURL, rcvHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
