package sqs

import (
	"sync"
	"time"

	"github.com/HomesNZ/go-common/sns"
)

type Event struct {
	Type    string    `json:"type"`
	Created time.Time `json:"created"`
}

func NewEvent(eventType string) Event {
	return Event{
		Type:    eventType,
		Created: time.Now(),
	}
}

func (ev Event) EventType() string {
	return ev.Type
}

type EventTyper interface {
	EventType() string
}

type EventSender interface {
	Send(ev EventTyper) error
}

type eventSender struct {
	mu     sync.Mutex
	topics map[string]*sns.Topic
}

func NewEventSender() EventSender {
	return &eventSender{
		topics: make(map[string]*sns.Topic),
	}
}

func (e *eventSender) Send(ev EventTyper) error {
	topic, err := e.initTopic(ev.EventType())
	if err != nil {
		return err
	}
	_, err = topic.PushMessage(ev)
	return err
}

func (e *eventSender) initTopic(name string) (*sns.Topic, error) {
	topic, ok := e.getTopic(name)
	if ok {
		return topic, nil
	}
	topic, err := sns.NewTopic(name)
	if err != nil {
		return nil, err
	}
	e.setTopic(name, topic)
	return topic, nil
}

func (e *eventSender) getTopic(name string) (*sns.Topic, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	topic, ok := e.topics[name]
	return topic, ok
}

func (e *eventSender) setTopic(name string, topic *sns.Topic) {
	e.mu.Lock()
	e.topics[name] = topic
	e.mu.Unlock()
}
