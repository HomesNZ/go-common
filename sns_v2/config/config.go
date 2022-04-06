package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config struct {
	Region           string // - is aws SQS region
	MessageStructure string
	Env              string
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Region, validation.Required.Error("AWS_REGION was not provided")),
	)
}
