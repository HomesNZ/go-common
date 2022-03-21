package config

import "github.com/HomesNZ/go-common/env"

func NewFromEnv() (*Config, error) {
	accessKeyID := env.GetString("AWS_ACCESS_KEY_ID", "")
	secretAccessKey := env.GetString("AWS_SECRET_ACCESS_KEY", "")
	region := env.GetString("AWS_SQS_REGION", "")
	messageStructure := env.GetString("AWS_SNS_MESSAGE_STRUCTURE", "json")
	suffix := env.Env()
	if suffix == "" {
		suffix = "development"
	}

	cfg := &Config{
		AccessKeyID:      accessKeyID,
		SecretAccessKey:  secretAccessKey,
		Region:           region,
		MessageStructure: messageStructure,
		Env:              suffix,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
