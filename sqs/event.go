package sqs

import (
	"time"
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