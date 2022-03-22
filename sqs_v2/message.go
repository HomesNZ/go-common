package sqs_v2

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"time"
)

type Message struct {
	Type             string
	MessageID        string `json:"MessageId"`
	TopicArn         string
	Message          string
	Timestamp        time.Time
	SignatureVersion string
	Signature        string
	SigningCertURL   string
	UnsubscribeURL   string

	sqsMessage types.Message
}

func newMessage(sqsMessage types.Message) (Message, error) {
	m := Message{
		sqsMessage: sqsMessage,
	}
	err := json.Unmarshal([]byte(*sqsMessage.Body), &m)
	return m, err
}
