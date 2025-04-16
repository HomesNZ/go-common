package sns_v2

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/HomesNZ/go-common/sns_v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/google/uuid"
)

type Service interface {
	Send(ctx context.Context, eventType string, message interface{}) error
	SendBatch(ctx context.Context, eventType string, messages []interface{}) error
}

type TopicArn *string

type service struct {
	conn   *sns.Client
	config *config.Config
	mu     sync.RWMutex
	topics map[string]TopicArn
}

func (s *service) Send(ctx context.Context, eventType string, message interface{}) error {
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
		MessageStructure: &s.config.MessageStructure,
		TopicArn:         topicArn,
		Message:          &m,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *service) SendBatch(ctx context.Context, eventType string, messages []interface{}) error {
	topic, err := s.topic(ctx, eventType)
	if err != nil {
		return err
	}

	if len(messages) == 0 {
		return nil
	}

	var publishBatchRequestEntries []types.PublishBatchRequestEntry
	for _, m := range messages {
		messageObjBytes, err := json.Marshal(m)
		if err != nil {
			return err
		}

		messageWrapper := Message{string(messageObjBytes)}
		messageBytes, err := json.Marshal(messageWrapper)
		if err != nil {
			return err
		}
		message := string(messageBytes)
		id := uuid.NewString()
		publishBatchRequestEntries = append(publishBatchRequestEntries, types.PublishBatchRequestEntry{
			Id:               &id,
			Message:          &message,
			MessageStructure: &s.config.MessageStructure,
		})
	}

	_, err = s.conn.PublishBatch(
		ctx,
		&sns.PublishBatchInput{
			TopicArn:                   topic,
			PublishBatchRequestEntries: publishBatchRequestEntries,
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) name(name string) *string {
	n := name + "_" + s.config.Env
	return &n
}

func (s *service) topic(ctx context.Context, name string) (*string, error) {
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

func (s *service) getTopic(name string) (*string, bool) {
	s.mu.RLock()
	topic, ok := s.topics[name]
	s.mu.RUnlock()
	return topic, ok
}

func (s *service) setTopic(name string, topic TopicArn) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.topics[name] = topic
}
