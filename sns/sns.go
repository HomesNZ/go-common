package sns

import (
	"context"
	"encoding/json"
	"github.com/HomesNZ/go-common/sns/config"
	"sync"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type Service interface {
	Send(ctx context.Context, eventType string, message interface{}) error
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

func (s *service) name(name string) string {
	return name + "_" + s.config.Env
}

func (s *service) topic(ctx context.Context, name string) (*string, error) {
	if topic, ok := s.getTopic(name); ok {
		return topic, nil
	}

	input := &sns.CreateTopicInput{
		Name: &name,
	}
	output, err := s.conn.CreateTopic(ctx, input)
	if err != nil {
		return nil, err
	}
	s.setTopic(name, output.TopicArn)
	return output.TopicArn, nil
}

func (s *service) getTopic(name string) (*string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	topic, ok := s.topics[name]
	return topic, ok
}

func (s *service) setTopic(name string, topic TopicArn) {
	s.mu.RLock()
	s.topics[name] = topic
	s.mu.RUnlock()
}
