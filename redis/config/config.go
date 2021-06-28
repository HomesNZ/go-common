package config

import (
	"fmt"
	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation"
)

type Config interface {
	Addr() string
}

type config struct {
	host string
	port string
}

func (c config) Addr() string {
	return fmt.Sprintf("%s:%s", c.host, c.port)
}

func New(host string, port string) (Config, error) {
	cfg := &config{
		host: host,
		port: port,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func NewFromEnv() (Config, error) {
	host := env.GetString("REDIS_HOST", "localhost")
	port := env.GetString("REDIS_PORT", "6379")

	cfg := &config{
		host: host,
		port: port,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c config) validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.host, validation.Required, validation.Required.Error("Redis host was not provided")),
		validation.Field(&c.port, validation.Required, validation.Required.Error("Redis port was not provided")),
	)
}
