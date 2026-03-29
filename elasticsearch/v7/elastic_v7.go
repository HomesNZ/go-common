package v7

import (
	"strings"

	"github.com/HomesNZ/env"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/pkg/errors"
	"github.com/sha1sum/aws_signing_client"

	elastic "github.com/olivere/elastic/v7"
)

// Gets a connection to AWS OpenSearch - depends on ENV:
//
//	ELASTICSEARCH_URLS
//	--testing:
//	+ ELASTICSEARCH_INSECURE=true
//	--way of AWS:
//	+ AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
//	+ AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
//	+ AWS_REGION (optional)
func Conn() (*elastic.Client, error) {
	url_raw := env.GetString("ELASTICSEARCH_URLS", "")
	if url_raw == "" {
		return nil, errors.New("missing ELASTICSEARCH_URLS")
	}
	urls := strings.Split(url_raw, ";")

	params := []elastic.ClientOptionFunc{
		elastic.SetURL(urls...),
		elastic.SetMaxRetries(10),
		// https://github.com/olivere/elastic/wiki/Using-with-AWS-Elasticsearch-Service
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetGzip(false),
	}
	// should be used only for testing / localhost
	insecureReq := env.GetString("ELASTICSEARCH_INSECURE", "false") == "true"
	if !insecureReq {
		signer := v4.NewSigner(credentials.NewEnvCredentials())
		region := env.GetString("AWS_REGION", "ap-southeast-2")
		awsClient, err := aws_signing_client.New(signer, nil, "es", region)
		if err != nil {
			return nil, errors.Wrap(err, "error creating aws signer")
		}
		params = append(params, elastic.SetHttpClient(awsClient))
	}

	client, err := elastic.NewClient(params...)
	if err != nil {
		return nil, errors.Wrap(err, "error creating client")
	}
	return client, nil
}
