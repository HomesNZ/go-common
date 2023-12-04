package middleware_test

import (
	"context"
	"os"
	"testing"

	lambda "github.com/HomesNZ/go-common/lambda/middleware"
	"github.com/aws/aws-lambda-go/events"
)

var requestHeaderOrigin = "http://localhost:3000"

func TestLambdaMiddleware(t *testing.T) {

	// Should return 200 when origin is allowed and the origin is set in the ENV and matched
	// and it should not override the existing headers
	os.Setenv("REQUEST_HEADER_ORIGIN", "http://localhost:3000;http://localhost:6000")
	res, err := runTest()
	if err != nil {
		t.Error(err)
	}
	if res.StatusCode != 200 {
		t.Error("StatusCode is not 200")
	} else if res.Body != "test_body" {
		t.Error("Body is not test_body")
	} else if res.Headers["Access-Control-Allow-Origin"] != requestHeaderOrigin {
		t.Errorf("Access-Control-Allow-Origin is not %s", requestHeaderOrigin)
	} else if res.Headers["Access-Control-Allow-Credentials"] != "true" {
		t.Errorf("Access-Control-Allow-Credentials is not true")
	} else if res.Headers["Content-Type"] != "application/json" {
		t.Errorf("Content-Type is not application/json")
	}

	// Should return error when origin is not allowed, the origin is not set in the ENV
	lambda.ResetRequestHeaderOriginEnv()
	os.Setenv("REQUEST_HEADER_ORIGIN", "")
	res, err = runTest()
	if err == nil {
		t.Error("error should exist when RequestHeaderOriginEnv is not set")
	}

	// Should return error when origin is not allowed, the origin is not matched in the ENV
	lambda.ResetRequestHeaderOriginEnv()
	os.Setenv("REQUEST_HEADER_ORIGIN", "http://localhost:6000")
	res, err = runTest()
	if err == nil {
		t.Error(err)
	}
}

func runTest() (events.APIGatewayProxyResponse, error) {
	return lambda.AddResponseHeaders(testLambdaHandler)(context.Background(), events.APIGatewayProxyRequest{
		Headers: map[string]string{
			"Content-Type": "application/json",
			"origin":       requestHeaderOrigin,
		},
	})
}

func testLambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "test_body",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
