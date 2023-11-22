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
	mapRequestHeaderOrigin = make(map[string]struct{})
)

type LambdaHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func AddResponseHeaders(next LambdaHandler) LambdaHandler {
	requestHeaderOriginEnv = env.GetString("REQUEST_HEADER_ORIGIN", "")
	if requestHeaderOriginEnv != "" {
		origins := strings.Split(requestHeaderOriginEnv, ";")
		for _, value := range origins {
			origin := strings.TrimSpace(value)
			mapRequestHeaderOrigin[origin] = struct{}{}
		}
	}

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		if requestHeaderOriginEnv == "" {
			return events.APIGatewayProxyResponse{}, errors.New("REQUEST_HEADER_ORIGIN is empty")
		}

		res, err := next(ctx, request)
		if err != nil {
			return res, err
		}

		requestHeaderOrigin := request.Headers["origin"]
		if _, ok := mapRequestHeaderOrigin[requestHeaderOrigin]; ok {
			if res.Headers == nil {
				res.Headers = make(map[string]string)
			}
			res.Headers["Access-Control-Allow-Origin"] = requestHeaderOrigin
			res.Headers["Access-Control-Allow-Credentials"] = "true"
		}

		return res, nil
	}
}

// ResetRequestHeaderOriginEnv resets the RequestHeaderOriginEnv variable, it is used for testing
func ResetRequestHeaderOriginEnv() {
	requestHeaderOriginEnv = ""
	mapRequestHeaderOrigin = make(map[string]struct{})
}
