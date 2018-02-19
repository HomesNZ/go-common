package elasticsearch

import (
	"net/http"

	"github.com/HomesNZ/elastic"
	awsauth "github.com/smartystreets/go-aws-auth"
)

func New(cfg ...elastic.ClientOptionFunc) (*elastic.Client, error) {
	// Set some defaults
	options := []elastic.ClientOptionFunc{
		elastic.SetSniff(false), // causes issues within AWS, so off by default
	}

	return elastic.NewClient(append(options, cfg...)...)
}

func AWSAccessKey(accessKeyID, secretAccessKey string) elastic.ClientOptionFunc {
	return elastic.SetPrepareRequest(func(req *http.Request) {
		awsauth.Sign(req, awsauth.Credentials{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
		})
	})
}

func AWSSecurityToken(securityToken string) elastic.ClientOptionFunc {
	return elastic.SetPrepareRequest(func(req *http.Request) {
		awsauth.Sign(req, awsauth.Credentials{
			SecurityToken: securityToken,
		})
	})
}
