package sns

import (
	"sync"

	"github.com/HomesNZ/go-common/env"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sirupsen/logrus"
)

var (
	conn     *sns.SNS
	initOnce = sync.Once{}

	contextLogger = logrus.WithField("package", "sns")
)

func snsEnabled() bool {
	enabled := env.GetBoolOrFalse("ENABLE_SNS_EVENTS")
	if !enabled {
		contextLogger.Info("SNS method called but ENABLE_SNS_EVENTS=false")
	}
	return enabled
}

func initConn() {
	if snsEnabled() {
		conn = sns.New(session.New(), &aws.Config{
			Region:      aws.String(env.GetString("SNS_REGION", "")),
			Credentials: credentials.NewEnvCredentials(),
		})
	}
}

// Conn returns the SNS connection
func Conn() *sns.SNS {
	initOnce.Do(initConn)
	return conn
}
