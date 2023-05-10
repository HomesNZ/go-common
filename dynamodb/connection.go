package dynamodb

import (
	"sync"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/guregu/dynamo"
)

var (
	conn     *dynamo.DB
	initOnce = sync.Once{}
)

func initConn() {
	conn = dynamo.New(session.New(), &aws.Config{
		Region:      aws.String(env.GetString("DYNAMODB_REGION", "")),
		Credentials: credentials.NewEnvCredentials(),
	})
}

// Conn returns the connection to DynamoDB
func Conn() *dynamo.DB {
	initOnce.Do(initConn)
	return conn
}
