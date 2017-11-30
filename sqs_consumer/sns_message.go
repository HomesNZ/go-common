package sqsConsumer

import (
	"encoding/json"
	"time"

	"github.com/goamz/goamz/sqs"
)

// SNSMessageHandler is the same as MessageHandler except it converts an SQS
// message to an SNS message format before sending to the handler.
type SNSMessageHandler func(message SNSMessage) bool

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

	sqsMessage *sqs.Message
}

func newSNSMessage(sqsMessage *sqs.Message) (SNSMessage, error) {
	m := SNSMessage{
		sqsMessage: sqsMessage,
	}
	err := json.Unmarshal([]byte(sqsMessage.Body), &m)
	return m, err
}

// ChangeMessageVisibility sets the visibility timeout for a given message.
// http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html
func (m SNSMessage) ChangeMessageVisibility(consumer Consumer, d time.Duration) error {
	_, err := consumer.queue.ChangeMessageVisibility(
		m.sqsMessage,
		int(d.Seconds()),
	)
	return err
}
