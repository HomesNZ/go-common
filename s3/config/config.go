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
type config struct {
	bucketName      string
	acl             string
	region          string
	endpoint        string
	cloudfrontURL   string
	accessKeyID     string
	secretAccessKey string
}

type Config interface {
	Validate() error
	BucketName() string
	ACL() string
	CloudfrontURL() string
	Region() string
	Endpoint() string
	AccessKeyID() string
	SecretAccessKey() string
}

func NewFromEnv() (Config, error) {
	cfg := &config{
		acl:             env.GetString("AWS_S3_ACL", "private"),
		region:          env.GetString("AWS_S3_REGION", "ap-southeast-2"),
		endpoint:        env.GetString("AWS_S3_ENDPOINT", "s3-ap-southeast-2.amazonaws.com"),
		bucketName:      env.GetString("AWS_S3_BUCKET", ""),
		cloudfrontURL:   env.GetString("CDN_URL", ""),
		accessKeyID:     env.GetString("AWS_ACCESS_KEY_ID", ""),
		secretAccessKey: env.GetString("AWS_SECRET_ACCESS_KEY", ""),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *config) BucketName() string {
	return c.bucketName
}

func (c *config) ACL() string {
	return c.acl
}

func (c *config) CloudfrontURL() string {
	return c.cloudfrontURL
}

func (c *config) Region() string {
	return c.region
}

func (c *config) Endpoint() string {
	return c.endpoint
}

func (c *config) AccessKeyID() string {
	return c.accessKeyID
}

func (c *config) SecretAccessKey() string {
	return c.secretAccessKey
}

func (c *config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.accessKeyID, validation.Required, validation.Required.Error("AWS access key was not provided")),
		validation.Field(&c.secretAccessKey, validation.Required, validation.Required.Error("AWS secret access key was not provided")),
		validation.Field(&c.bucketName, validation.Required, validation.Required.Error("Bucket name was not provided")),
	)
}
