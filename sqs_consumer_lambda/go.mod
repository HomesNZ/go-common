module github.com/HomesNZ/go-common/sqs_consumer_lambda

go 1.16

require (
	github.com/HomesNZ/events v0.0.0-20210526041501-6acb2a727cf4
	github.com/HomesNZ/go-common/redis v0.0.0-20210429042447-401b06a429ea
	github.com/aws/aws-lambda-go v1.24.0
	github.com/go-redsync/redsync v1.4.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
)
