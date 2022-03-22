package sns_v2

import (
	"context"
	"github.com/HomesNZ/go-common/sns_v2/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewFromEnv(ctx context.Context) (Service, error) {

	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}
	client := sns.NewFromConfig(aws.Config{
		Region:      config.Region,
		Credentials: credentials.NewStaticCredentialsProvider(config.AccessKeyID, config.SecretAccessKey, ""),
	})

	return &service{conn: client, config: config}, nil
}