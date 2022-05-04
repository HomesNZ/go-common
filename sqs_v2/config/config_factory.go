package config

import "github.com/HomesNZ/go-common/env"

func New(region, queueName string) (*Config, error) {
	cfg := &Config{
		QueueName: queueName,
		Region:    region,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromEnv() (*Config, error) {
	region := env.GetString("AWS_SQS_REGION", "")
	queueName := env.GetString("AWS_SQS_QUEUE", "")
	maxMsg := env.GetInt("AWS_SQS_MAX_MESSAGES", 1)
	maxWorker := env.GetInt("AWS_SQS_MAX_WORKERS", 1)

	cfg := &Config{
		QueueName: queueName,
		Region:    region,
		MaxMsg:    int32(maxMsg),
		MaxWorker: maxWorker,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
