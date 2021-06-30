package config

import "github.com/HomesNZ/go-common/env"

func New(accessKeyID, secretAccessKey, region, queueName string) (Config, error) {
	cfg := &config{
		queueName:       queueName,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		region:          region,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromEnv() (Config, error) {
	accessKeyID := env.GetString("AWS_ACCESS_KEY_ID", "")
	secretAccessKey := env.GetString("AWS_SECRET_ACCESS_KEY", "")
	region := env.GetString("AWS_SQS_REGION", "")
	queueName := env.GetString("AWS_SQS_QUEUE", "")

	cfg := &config{
		queueName:       queueName,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
		region:          region,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
