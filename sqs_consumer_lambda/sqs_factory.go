package sqsConsumerLambda

import (
	"github.com/pkg/errors"
)

// New for AWS Lambda
func New(handlers map[string]SNSMessageHandler) (Consumer, error) {
	if len(handlers) == 0 {
		return nil, errors.New("no handlers provided")
	}

	router := newRouter()
	for event, h := range handlers {
		router.AddRoute(event, h)
	}

	handler := Handler{
		Router: router,
	}


	return &consumer{
		handler:           SNSMessageHandler(handler.HandleMessage),
		waitForCompletion: true,
	}, nil
}

