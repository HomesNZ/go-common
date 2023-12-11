package sns_v2

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/HomesNZ/go-common/trace"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/HomesNZ/go-common/sns_v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

const (
	attrHomesTrace = "homes_trace"
)

type TopicArn *string

type Service struct {
	conn   *sns.Client
	config *config.Config
	mu     sync.RWMutex
	topics map[string]TopicArn
}

func (s *Service) Send(ctx context.Context, eventType string, message interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}

	// If the context already has a trace, use it. Otherwise, create a new trace.
	// we don't need to link the trace to the parent trace because this responsibility of the receiver not the sender
	eventTrace := trace.FromCtx(ctx)
	if eventTrace.IsEmpty() {
		eventTrace = trace.New()
	}

	attrs := map[string]types.MessageAttributeValue{
		attrHomesTrace: {
			DataType:    aws.String("String"),
			StringValue: aws.String(eventTrace.ToJSON()),
		},
	}

	topicArn, err := s.topic(ctx, eventType)
	if err != nil {
		return err
	}
	messageObjBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	messageWrapper := Message{string(messageObjBytes)}
	messageBytes, err := json.Marshal(messageWrapper)
	if err != nil {
		return err
	}
	m := string(messageBytes)
	_, err = s.conn.Publish(ctx, &sns.PublishInput{
		MessageStructure:  &s.config.MessageStructure,
		TopicArn:          topicArn,
		Message:           &m,
		MessageAttributes: attrs,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) name(name string) *string {
	n := name + "_" + s.config.Env
	return &n
}

func (s *Service) topic(ctx context.Context, name string) (*string, error) {
	if topic, ok := s.getTopic(name); ok {
		return topic, nil
	}

	input := &sns.CreateTopicInput{
		Name: s.name(name),
	}
	output, err := s.conn.CreateTopic(ctx, input)
	if err != nil {
		return nil, err
	}
	s.setTopic(name, output.TopicArn)
	return output.TopicArn, nil
}

func (s *Service) getTopic(name string) (*string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	topic, ok := s.topics[name]
	return topic, ok
}

func (s *Service) setTopic(name string, topic TopicArn) {
	s.mu.RLock()
	s.topics[name] = topic
	s.mu.RUnlock()
}
