package sqsConsumerLambda

import (
	"github.com/aws/aws-lambda-go/events"
)

type SQSMessage events.SQSMessage
type SQSMessageHandler func(message SQSMessage) (bool, error)
