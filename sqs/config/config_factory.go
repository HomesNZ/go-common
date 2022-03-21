package config

import "github.com/HomesNZ/go-common/env"

func New(accessKeyID, secretAccessKey, region, queueName string) (*Config, error) {
	cfg := &Config{
		QueueName:       queueName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromEnv() (*Config, error) {
	accessKeyID := env.GetString("AWS_ACCESS_KEY_ID", "")
	secretAccessKey := env.GetString("AWS_SECRET_ACCESS_KEY", "")
	region := env.GetString("AWS_SQS_REGION", "")
	queueName := env.GetString("AWS_SQS_QUEUE", "")

	cfg := &Config{
		QueueName:       queueName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		Region:          region,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
