package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/pkg/errors"
)

type logger interface {
	Error(ctx context.Context, msg string, args ...any)
	Info(ctx context.Context, msg string, args ...any)
}

func ErrorHandler(log logger, internalError error) func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Error(context.TODO(), internalError.Error())
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return SendError(http.StatusInternalServerError, internalError)
	}
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

func SendError(status int, err error) (events.APIGatewayProxyResponse, error) {
	var statusCode int
	if status == 0 || status >= 600 || status < 100 {
		statusCode = http.StatusInternalServerError
	} else if status >= 100 && status < 600 {
		statusCode = status
	}

	responseBody := Response{
		Status:  statusCode,
		Message: err.Error(),
	}
	body, _ := json.Marshal(responseBody)

	return events.APIGatewayProxyResponse{
		Headers:    map[string]string{"Content-Type": "application/json"},
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}

func Send(request events.APIGatewayProxyRequest, statusCode int, data any) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return SendError(http.StatusInternalServerError, errors.Wrap(err, "failed to marshal response body"))
	}
	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      request.Headers["origin"],
			"Access-Control-Allow-Credentials": "true",
		},
		StatusCode: statusCode,
		Body:       string(body),
	}, nil
}
