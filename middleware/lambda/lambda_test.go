package middleware_lambda_test

import (
	"context"
	"testing"

	lambda "github.com/HomesNZ/go-common/middleware/lambda"
	"github.com/aws/aws-lambda-go/events"
)

func TestService(t *testing.T) {
	lambda.RequestHeaderOriginEnv = "http://localhost:3000"
	res, err := runTest()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("StatusCode is not 200")
	} else if res.Body != "test_body" {
		t.Error("Body is not test_body")
	} else if res.Headers["Access-Control-Allow-Origin"] != lambda.RequestHeaderOriginEnv {
		t.Errorf("Access-Control-Allow-Origin is not %s", lambda.RequestHeaderOriginEnv)
	} else if res.Headers["Access-Control-Allow-Credentials"] != "true" {
		t.Errorf("Access-Control-Allow-Credentials is not true")
	}

	lambda.RequestHeaderOriginEnv = ""
	res, err = runTest()
	if err == nil {
		t.Error("error should exist when RequestHeaderOriginEnv is not set")
	}
}

func runTest() (events.APIGatewayProxyResponse, error) {
	return lambda.AddResponseHeaders(testLambdaHandler)(context.Background(), events.APIGatewayProxyRequest{})
}

func testLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "test_body",
	}, nil
}
