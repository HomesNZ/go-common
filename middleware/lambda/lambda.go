package middleware_lambda

import (
	"context"
	"errors"
	"strings"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-lambda-go/events"
)

var (
	requestHeaderOriginEnv string
	mapRequestHeaderOrigin = make(map[string]bool)
)

func init() {
	requestHeaderOriginEnv = env.GetString("REQUEST_HEADER_ORIGIN", "")
	for _, value := range strings.Split(requestHeaderOriginEnv, ";") {
		mapRequestHeaderOrigin[value] = true
	}
}

type LambdaHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func AddResponseHeaders(next LambdaHandler) LambdaHandler {
	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if requestHeaderOriginEnv == "" {
			return events.APIGatewayProxyResponse{}, errors.New("REQUEST_HEADER_ORIGIN does not exist in env")
		}

		res, err := next(ctx, request)
		if err != nil {
			return res, err
		}

		requestHeaderOrigin := request.Headers["origin"]
		if _, ok := mapRequestHeaderOrigin[requestHeaderOrigin]; ok {
			requestHeaderOriginEnv = requestHeaderOrigin
		}

		res.Headers = map[string]string{
			"Access-Control-Allow-Origin":      requestHeaderOriginEnv,
			"Access-Control-Allow-Credentials": "true",
		}
		return res, nil
	}
}
