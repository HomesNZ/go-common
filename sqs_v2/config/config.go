package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	QueueName            string // - is aws SQS queue name
	Region               string // - is aws SQS region
	MaxMessageHandleTime int    // - is the maximum time in seconds that a message is allowed to be handled
	MaxMsg               int32
	MaxWorker            int
	MaxHandlers          int
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.QueueName, validation.Required.Error("AWS_SQS_QUEUE was not provided")),
		validation.Field(&c.Region, validation.Required.Error("AWS_REGION was not provided")),
	)
}
