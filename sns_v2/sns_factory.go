package sns_v2

import (
	"context"

	"github.com/HomesNZ/go-common/sns_v2/config"
	awsCfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func NewFromEnv(ctx context.Context) (Service, error) {

	config, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	cfg, err := awsCfg.LoadDefaultConfig(ctx, awsCfg.WithRegion(config.Region))
	if err != nil {
		return nil, err
	}

	client := sns.NewFromConfig(cfg)

	return &service{conn: client, config: config}, nil
}
