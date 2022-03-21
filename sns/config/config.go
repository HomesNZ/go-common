package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	AccessKeyID      string // - is aws access key id
	SecretAccessKey  string // - is aws secret access key
	Region           string // - is aws SQS region
	MessageStructure string
	Env              string
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.AccessKeyID, validation.Required.Error("AWS_ACCESS_KEY_ID was not provided")),
		validation.Field(&c.SecretAccessKey, validation.Required.Error("AWS_SECRET_ACCESS_KEY was not provided")),
		validation.Field(&c.Region, validation.Required.Error("AWS_REGION was not provided")),
	)
}
