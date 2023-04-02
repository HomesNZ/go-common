package config

import (
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	APIKey string
	Stage  string
}

func NewFromEnv() *Config {
	cfg := &Config{
		APIKey: env.GetString("BUGSNAG_API_KEY", ""),
		Stage:  env.Env(),
	}

	return cfg
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.APIKey, validation.Required, validation.Required.Error("BUGSNAG_API_KEY was not specified")),
	)
}
