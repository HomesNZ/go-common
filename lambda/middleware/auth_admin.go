package middleware

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"net/http"
	"strings"

	"github.com/HomesNZ/go-secret/auth"
	"github.com/aws/aws-lambda-go/events"
)

var jwkAuth *auth.Auth
var (
	errAuthInvalid = errors.New("auth token invalid")
)

func AuthWithAdminRole(authURL string, next LambdaHandler) (LambdaHandler, error) {
	var err error
	jwkAuth, err = auth.New(
		auth.JWKS(authURL),
	)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		httpRequest, err := http.NewRequest(
			strings.ToUpper(req.HTTPMethod),
			req.Path,
			bytes.NewReader([]byte(req.Body)),
		)

		if req.MultiValueHeaders != nil {
			for k, values := range req.MultiValueHeaders {
				for _, value := range values {
					httpRequest.Header.Add(k, value)
				}
			}
		} else {
			for h := range req.Headers {
				httpRequest.Header.Add(h, req.Headers[h])
			}
		}

		httpRequest.RequestURI = httpRequest.URL.RequestURI()

		token, err := jwkAuth.Authenticate(httpRequest)
		if err != nil {
			return events.APIGatewayProxyResponse{}, errAuthInvalid
		}

		if !token.Valid || token.User == nil {
			return events.APIGatewayProxyResponse{}, errAuthInvalid
		}

		if !token.User.HasRole(auth.RoleAdmin) {
			return events.APIGatewayProxyResponse{}, errAuthInvalid
		}

		return next(ctx, req)
	}, nil
}
