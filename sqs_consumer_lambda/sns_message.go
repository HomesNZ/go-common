package sqsConsumerLambda

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

// copy from AWS sdk
type MessageAttributeValue struct {

	// Amazon SNS supports the following logical data types: String, String.Array,
	// Number, and Binary. For more information, see Message Attribute Data Types
	// (https://docs.aws.amazon.com/sns/latest/dg/SNSMessageAttributes.html#SNSMessageAttributes.DataTypes).
	//
	// This member is required.
	DataType *string

	// Binary type attributes can store any binary data, for example, compressed data,
	// encrypted data, or images.
	BinaryValue []byte

	// Strings are Unicode with UTF8 binary encoding. For a list of code values, see
	// ASCII Printable Characters
	// (https://en.wikipedia.org/wiki/ASCII#ASCII_printable_characters).
	StringValue *string
	// contains filtered or unexported fields
}

// SNSMessageHandler is the same as MessageHandler except it converts an SQS
// message to an SNS message format before sending to the handler.

type SNSMessageHandler func(message SNSMessage) (bool, error)

// SNSMessage is a data struct matching the output from a message pushed through
// SQS from SNS.
type SNSMessage struct {
	Type              string
	MessageID         string `json:"MessageId"`
	TopicArn          string
	Message           string
	Timestamp         time.Time
	SignatureVersion  string
	Signature         string
	SigningCertURL    string
	UnsubscribeURL    string
	MessageAttributes map[string]MessageAttributeValue

	sqsMessage *events.SQSMessage
}

func newSNSMessage(sqsMessage *events.SQSMessage) (SNSMessage, error) {
	m := SNSMessage{
		sqsMessage: sqsMessage,
	}
	err := json.Unmarshal([]byte(sqsMessage.Body), &m)
	return m, err
}
