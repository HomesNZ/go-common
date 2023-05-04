package middleware_lambda

import (
	"context"
	"errors"
	"strings"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-lambda-go/events"
)

var (
	RequestHeaderOriginEnv string
	mapRequestHeaderOrigin = make(map[string]bool)
)

func init() {
	RequestHeaderOriginEnv = env.GetString("REQUEST_HEADER_ORIGIN", "")
	for _, value := range strings.Split(RequestHeaderOriginEnv, ";") {
		mapRequestHeaderOrigin[value] = true
	}
}

type LambdaHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func AddResponseHeaders(next LambdaHandler) LambdaHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if RequestHeaderOriginEnv == "" {
			return events.APIGatewayProxyResponse{}, errors.New("REQUEST_HEADER_ORIGIN is empty")
		}

		res, err := next(ctx, request)
		if err != nil {
			return res, err
		}

		requestHeaderOrigin := request.Headers["origin"]
		if _, ok := mapRequestHeaderOrigin[requestHeaderOrigin]; ok {
			RequestHeaderOriginEnv = requestHeaderOrigin
		}

		res.Headers = map[string]string{
			"Access-Control-Allow-Origin":      RequestHeaderOriginEnv,
			"Access-Control-Allow-Credentials": "true",
		}
		return res, nil
	}
}
