package config

import (
	"time"

	"github.com/HomesNZ/go-common/env"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Config struct {
	ServiceName       string
	Host              string
	User              string
	Name              string
	Password          string
	MaxConns          int // max number of connections in the pool. will define
	Port              int
	SearchPath        string
	HealthCheckPeriod time.Duration // seconds - how often to check health of the connection
	MaxConnIdleTime   time.Duration // minutes - how long connection can be idle before it'll be closed
	PingBeforeUse     bool          // it'll be used to check connection before use, if true and connection is not alive, it'll be reconnected
}

func NewFromEnv() *Config {
	healthCheckPeriod := time.Duration(env.GetInt("DB_HEALTH_CHECK_PERIOD", 30)) * time.Second
	maxConnIdleTime := time.Duration(env.GetInt("DB_MAX_CONN_IDLE_TIME", 5)) * time.Minute
	pingBeforeUse := env.GetBool("DB_PING_BEFORE_USE", true)

	cfg := &Config{
		ServiceName:       env.GetString("SERVICE_NAME", ""),
		Host:              env.GetString("DB_HOST", "localhost"),
		User:              env.GetString("DB_USER", "postgres"),
		Name:              env.GetString("DB_NAME", ""),
		Password:          env.GetString("DB_PASSWORD", ""),
		MaxConns:          env.GetInt("DB_MAX_CONNECT", 3),
		Port:              env.GetInt("DB_PORT", 5432),
		SearchPath:        env.GetString("DB_SEARCH_PATH", ""),
		HealthCheckPeriod: healthCheckPeriod,
		MaxConnIdleTime:   maxConnIdleTime,
		PingBeforeUse:     pingBeforeUse,
	}

	return cfg
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.ServiceName, validation.Required, validation.Required.Error("SERVICE_NAME was not specified in env")),
		validation.Field(&c.Name, validation.Required, validation.Required.Error("DB_NAME was not specified in env")),
	)
}
