package redis

import (
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation"
)

type DeclareConfig func(config *Config)

type Config struct {
	Host string
	Port string
}

func ConfigFromEnv() (*Config, error) {
	host := env.GetString("REDIS_HOST", "localhost")
	port := env.GetString("REDIS_PORT", "6379")

	cfg := &Config{
		Host: host,
		Port: port,
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required, validation.Required.Error("Redis host was not provided")),
		validation.Field(&c.Port, validation.Required, validation.Required.Error("Redis port was not provided")),
	)
}
