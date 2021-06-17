package sqsConsumerLambda

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// SNSMessageHandler is the same as MessageHandler except it converts an SQS
// message to an SNS message format before sending to the handler.

type SNSMessageHandler func(message SNSMessage) (bool, error)

// SNSMessage is a data struct matching the output from a message pushed through
// SQS from SNS.
type SNSMessage struct {
	Type             string
	MessageID        string `json:"MessageId"`
	TopicArn         string
	Message          string
	Timestamp        time.Time
	SignatureVersion string
	Signature        string
	SigningCertURL   string
	UnsubscribeURL   string

	sqsMessage *events.SQSMessage
}

func newSNSMessage(sqsMessage *events.SQSMessage) (SNSMessage, error) {
	m := SNSMessage{
		sqsMessage: sqsMessage,
	}
	err := json.Unmarshal([]byte(sqsMessage.Body), &m)
	return m, err
}
