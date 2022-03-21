package sqs

import (
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSMessage types.Message
type MessageHandler func(message SQSMessage) (bool, error)
