package v7

import (
	"strings"

	"github.com/HomesNZ/env"
	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/sha1sum/aws_signing_client"
	"github.com/sirupsen/logrus"

	elastic "github.com/olivere/elastic/v7"
)

// Gets a connection to AWS OpenSearch - depends on ENV
//
//	ELASTICSEARCH_URLS
//	+ AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY
//	+ AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
func Conn() *elastic.Client {
	log := logrus.WithField("package", "elasticsearch")
	signer := v4.NewSigner(credentials.NewEnvCredentials())
	awsClient, err := aws_signing_client.New(signer, nil, "es", "ap-southeast-2")
	if err != nil {
		log.Errorln("error creating aws signer")
		panic(err)
	}

	urls := strings.Split(env.MustGetString("ELASTICSEARCH_URLS"), ";")
	client, err := elastic.NewClient(
		elastic.SetURL(urls...),
		elastic.SetMaxRetries(10),
		elastic.SetHttpClient(awsClient),
		// https://github.com/olivere/elastic/wiki/Using-with-AWS-Elasticsearch-Service
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
		elastic.SetGzip(false),
	)
	if err != nil {
		log.Errorln("error creating client")
		panic(err)
	}
	log.Infoln("using olivere/elastic/v7", urls)
	return client
}
