package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	QueueName       string // - is aws SQS queue name
	AccessKeyID     string // - is aws access key id
	SecretAccessKey string // - is aws secret access key
	Region          string // - is aws SQS region
	MaxMsg          int32
	MaxWorker       int
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.QueueName, validation.Required.Error("AWS_SQS_QUEUE was not provided")),
		validation.Field(&c.AccessKeyID, validation.Required.Error("AWS_ACCESS_KEY_ID was not provided")),
		validation.Field(&c.SecretAccessKey, validation.Required.Error("AWS_SECRET_ACCESS_KEY was not provided")),
		validation.Field(&c.Region, validation.Required.Error("AWS_REGION was not provided")),
	)
}
