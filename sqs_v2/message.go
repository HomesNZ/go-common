package sqs_v2

import (
	"encoding/json"
	"github.com/HomesNZ/go-common/trace"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"time"
)

const (
	attrHomesTrace = "homes_trace"
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
	Trace            trace.Trace

	sqsMessage types.Message
}

func newMessage(sqsMessage types.Message) (Message, error) {
	m := Message{
		sqsMessage: sqsMessage,
	}
	var msgTrace trace.Trace
	if sqsMessage.MessageAttributes != nil {
		if traceAttr, ok := sqsMessage.MessageAttributes[attrHomesTrace]; ok {
			msgTrace = trace.LinkFromJSON(traceAttr.StringValue) // set a new event id to the trace
		} else {
			msgTrace = trace.New()
		}
	}
	m.Trace = msgTrace

	err := json.Unmarshal([]byte(*sqsMessage.Body), &m)
	return m, err
}
