package sqsConsumerLambda

import (
	"encoding/json"

	"github.com/HomesNZ/events"
	"github.com/pkg/errors"
)

//type MessageHandler func(ctx context.Context, message sqs.SNSMessage) (bool, error)

type Router struct {
	routes map[string]SNSMessageHandler
}

func newRouter() *Router {
	return &Router{
		routes: map[string]SNSMessageHandler{},
	}
}

func (r *Router) AddRoute(route string, handler SNSMessageHandler) {
	r.routes[route] = handler
}

func (r *Router) Handle(message SNSMessage) (bool, error) {
	rawJSON := []byte(message.Message)
	genericEvent := &events.Event{}
	err := json.Unmarshal(rawJSON, genericEvent)
	if err != nil {
		return true, errors.Wrap(err, "unmarshal generic")
	}

	handler, ok := r.routes[genericEvent.Type]
	if !ok {
		return true, errors.New("unknown event type: " + genericEvent.Type)
	}

	return handler(message)
}
