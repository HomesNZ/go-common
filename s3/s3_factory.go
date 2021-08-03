package s3

import (
	"github.com/HomesNZ/go-common/s3/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
)

// New initializes a new S3. If cloudfrontURL is not nil, URLs returned from UploadAsset will return the assets URL on
// Cloudfront distibution, otherwise the raw S3 URL will be returned.
func New(cfg config.Config) (Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newService(cfg), nil
}

func NewFromEnv() (Service, error) {
	cfg, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}
	return newService(cfg), nil
}

func newService(cfg config.Config) Service {
	awsConfig := &aws.Config{
		Region:           aws.String(cfg.Region()),
		Endpoint:         aws.String(cfg.Endpoint()),
		S3ForcePathStyle: aws.Bool(true),
		Credentials: credentials.NewCredentials(&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     cfg.AccessKeyID(),
			SecretAccessKey: cfg.SecretAccessKey(),
		}}),
	}

	return &s3{
		client: awsS3.New(session.New(), awsConfig),
		config: cfg,
	}
}
