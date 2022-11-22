package config

import (
	"github.com/HomesNZ/env"
	validation "github.com/go-ozzo/ozzo-validation"
)

type PresignConfig struct {
	AccessKeyID     string
	SecretAccessKey string
	Region          string
}

func NewPresignConfigFromEnv() *PresignConfig {
	cfg := &PresignConfig{
		AccessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
		Region:          env.GetString("AWS_S3_REGION", "ap-southeast-2"),
	}

	return cfg
}

func (c *PresignConfig) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.AccessKeyID, validation.Required, validation.Required.Error("AWS access key was not provided")),
		validation.Field(&c.SecretAccessKey, validation.Required, validation.Required.Error("AWS secret access key was not provided")),
	)
}
