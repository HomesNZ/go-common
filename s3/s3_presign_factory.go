package s3

import (
	"context"

	"github.com/HomesNZ/go-common/s3/config"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	awsCred "github.com/aws/aws-sdk-go-v2/credentials"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func NewPresignService(ctx context.Context, cfg *config.PresignConfig) (PresignService, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return newPresignService(ctx, cfg)
}

func NewPresignServiceFromEnv(ctx context.Context) (PresignService, error) {
	cfg := config.NewPresignConfigFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return newPresignService(ctx, cfg)
}

func newPresignService(ctx context.Context, cfg *config.PresignConfig) (PresignService, error) {
	creds := awsCred.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")
	awsCfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithRegion(cfg.Region),
		awsConfig.WithCredentialsProvider(creds),
	)
	if err != nil {
		return nil, err
	}
	client := awsS3.NewFromConfig(awsCfg)
	presignClient := awsS3.NewPresignClient(client)
	return &presignS3{
		presignClient: presignClient,
		presignConfig: cfg,
	}, nil
}
