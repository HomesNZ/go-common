package sqsConsumer

import (
	"time"

	"github.com/goamz/goamz/sqs"
)

// MessageHandler is an anonymous function which is used to handle messages
// recieved from the SQS queue. It should handle errors internally and return a
// simple boolean to indicate if handling was successful.
type MessageHandler func(message SQSMessage) bool

type SQSMessage sqs.Message

// ChangeMessageVisibility sets the visibility timeout for a given message.
// http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-visibility-timeout.html
func (m SQSMessage) ChangeMessageVisibility(consumer Consumer, d time.Duration) error {
	sqsMessage := sqs.Message(m)
	_, err := consumer.queue.ChangeMessageVisibility(
		&sqsMessage,
		int(d.Seconds()),
	)
	return err
}
