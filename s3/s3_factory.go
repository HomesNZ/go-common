package s3

import (
	"context"

	"github.com/HomesNZ/go-common/s3/config"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCred "github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// New initializes a new S3. If cloudfrontURL is not nil, URLs returned from UploadAsset will return the assets URL on
// Cloudfront distibution, otherwise the raw S3 URL will be returned.
func New(ctx context.Context, cfg *config.Config) (Service, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newService(ctx, cfg)
}

func NewFromEnv(ctx context.Context) (Service, error) {
	cfg := config.NewFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return newService(ctx, cfg)
}

func newService(ctx context.Context, cfg *config.Config) (Service, error) {
	creds := awsCred.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, cfg.SessionToken)
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}
	client := awsS3.NewFromConfig(awsCfg)
	return &s3{
		client: client,
		config: cfg,
	}, nil
}
