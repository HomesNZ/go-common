package config

import "github.com/HomesNZ/go-common/env"

func NewFromEnv() (*Config, error) {
	region := env.GetString("AWS_SQS_REGION", "")
	messageStructure := env.GetString("AWS_SNS_MESSAGE_STRUCTURE", "json")
	suffix := env.Env()
	if suffix == "" {
		suffix = "development"
	}

	cfg := &Config{
		Region:           region,
		MessageStructure: messageStructure,
		Env:              suffix,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}
