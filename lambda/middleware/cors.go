package middleware

import (
	"context"
	"strings"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-lambda-go/events"
)

var (
	requestHeaderOriginEnv string
	mapRequestHeaderOrigin = make(map[string]struct{})
)

type LambdaHandler func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

func Cors(next LambdaHandler) LambdaHandler {
	requestHeaderOriginEnv = env.GetString("REQUEST_HEADER_ORIGIN", "")
	if requestHeaderOriginEnv != "" {
		origins := strings.Split(requestHeaderOriginEnv, ";")
		for _, value := range origins {
			origin := strings.TrimSpace(value)
			mapRequestHeaderOrigin[origin] = struct{}{}
		}
	}

	return func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		res, err := next(ctx, request)
		if err != nil {
			return res, err
		}

		if res.Headers == nil {
			res.Headers = make(map[string]string)
		}

		if len(mapRequestHeaderOrigin) > 0 {
			requestHeaderOrigin := request.Headers["origin"]
			if _, ok := mapRequestHeaderOrigin[requestHeaderOrigin]; ok {
				res.Headers["Access-Control-Allow-Origin"] = requestHeaderOrigin
				res.Headers["Access-Control-Allow-Credentials"] = "true"
			}
		} else {
			res.Headers["Access-Control-Allow-Origin"] = request.Headers["origin"]
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
