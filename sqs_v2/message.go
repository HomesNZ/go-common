package sqs_v2

import (
	"encoding/json"
	"time"

	"github.com/HomesNZ/go-common/trace"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	attrHomesTrace = "homes_trace"
)

type Message struct {
	Type              string
	MessageID         string `json:"MessageId"`
	TopicArn          string
	Message           string
	Timestamp         time.Time
	SignatureVersion  string
	Signature         string
	SigningCertURL    string
	UnsubscribeURL    string
	Trace             trace.Trace
	MessageAttributes map[string]TypeValue

	sqsMessage types.Message
}

type TypeValue struct {
	Type  string
	Value string
}

func newMessage(sqsMessage types.Message) (Message, error) {
	m := Message{
		sqsMessage: sqsMessage,
	}
	err := json.Unmarshal([]byte(*sqsMessage.Body), &m)
	var msgTrace trace.Trace
	if m.MessageAttributes != nil {
		if traceAttr, ok := m.MessageAttributes[attrHomesTrace]; ok {
			msgTrace = trace.LinkFromJSON(traceAttr.Value) // set a new event id to the trace
		} else {
			msgTrace = trace.New()
		}
	}
	m.Trace = msgTrace
	return m, err
}
