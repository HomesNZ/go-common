package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config interface {
	QueueName() string
	AccessKeyID() string
	SecretAccessKey() string
	Region() string
}

type config struct {
	queueName       string // - is aws SQS queue name
	accessKeyID     string // - is aws access key id
	secretAccessKey string // - is aws secret access key
	region          string // - is aws SQS region
}

func (c config) QueueName() string {
	return c.queueName
}

func (c config) AccessKeyID() string {
	return c.accessKeyID
}

func (c config) SecretAccessKey() string {
	return c.secretAccessKey
}

func (c config) Region() string {
	return c.region
}

func (c config) validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.queueName, validation.Required, validation.Required.Error("AWS_SQS_QUEUE was not provided")),
		validation.Field(&c.accessKeyID, validation.Required, validation.Required.Error("AWS_ACCESS_KEY_ID was not provided")),
		validation.Field(&c.secretAccessKey, validation.Required, validation.Required.Error("AWS_SECRET_ACCESS_KEY was not provided")),
		validation.Field(&c.region, validation.Required, validation.Required.Error("AWS_SQS_REGION was not provided")),
	)
}
