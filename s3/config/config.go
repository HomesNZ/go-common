package config

import (
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation"
)

// ACL is policy for S3 assets.
// Region is the default region that assets are uploaded to.
// Endpoint is the default endpoint to be used when uploading assets.
// BucketName is aws S3 bucket Name
// CloudfrontURL is CDN url
type Config struct {
	BucketName      string
	ACL             string
	Region          string
	Endpoint        string
	CloudfrontURL   string
	AccessKeyID     string
	SecretAccessKey string
}

func NewFromEnv() (*Config, error) {
	cfg := &Config{
		ACL:             env.GetString("AWS_S3_ACL", "private"),
		Region:          env.GetString("AWS_S3_REGION", "ap-southeast-2"),
		Endpoint:        env.GetString("AWS_S3_ENDPOINT", "s3-ap-southeast-2.amazonaws.com"),
		BucketName:      env.GetString("AWS_S3_BUCKET", ""),
		CloudfrontURL:   env.GetString("CDN_URL", ""),
		AccessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
		SecretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AccessKeyID, validation.Required, validation.Required.Error("AWS access key was not provided")),
		validation.Field(&c.SecretAccessKey, validation.Required, validation.Required.Error("AWS secret access key was not provided")),
		validation.Field(&c.BucketName, validation.Required, validation.Required.Error("Bucket name was not provided")),
	)
}
