package elasticsearch

import (
	"net/http"

	"github.com/HomesNZ/elastic"
	awsauth "github.com/smartystreets/go-aws-auth"
)

// AWSAccessKey configures the client to sign each outgoing request with AWS V4
// signatures, using an IAM access key ID / secret key.
func AWSAccessKey(accessKeyID, secretAccessKey string) elastic.ClientOptionFunc {
	return elastic.SetPrepareRequest(func(req *http.Request) {
		awsauth.Sign(req, awsauth.Credentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		})
	})
}

// AWSSecurityToken configures the client to use a security token to
// authenticate with AWS.
func AWSSecurityToken(securityToken string) elastic.ClientOptionFunc {
	return elastic.SetPrepareRequest(func(req *http.Request) {
		awsauth.Sign(req, awsauth.Credentials{
			SecurityToken: securityToken,
		})
	})
}
