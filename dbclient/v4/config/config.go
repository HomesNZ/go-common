package config

import (
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	ServiceName               string
	Host                      string
	User                      string
	Name                      string
	Password                  string
	MaxConns                  int
	Port                      int
	SearchPath                string
	StandardConformingStrings bool
}

func NewFromEnv() *Config {
	cfg := &Config{
		ServiceName:               env.GetString("SERVICE_NAME", ""),
		Host:                      env.GetString("DB_HOST", "localhost"),
		User:                      env.GetString("DB_USER", "postgres"),
		Name:                      env.GetString("DB_NAME", ""),
		Password:                  env.GetString("DB_PASSWORD", ""),
		MaxConns:                  env.GetInt("DB_MAX_CONNECT", 1),
		Port:                      env.GetInt("DB_PORT", 5432),
		SearchPath:                env.GetString("DB_SEARCH_PATH", ""),
		StandardConformingStrings: env.GetBoolOrFalse("DB_STANDARD_CONFORMING_STRINGS"),
	}

	return cfg
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ServiceName, validation.Required, validation.Required.Error("SERVICE_NAME was not specified in env")),
		validation.Field(&c.Name, validation.Required, validation.Required.Error("DB_NAME was not specified in env")),
	)
}
