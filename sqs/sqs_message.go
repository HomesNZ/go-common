package sqs

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	//"time"
)

type SQSMessage sqs.Message
type MessageHandler func(message SQSMessage) (bool, error)
