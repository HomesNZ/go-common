package dbclient

import (
	"context"
	"fmt"

	"github.com/HomesNZ/go-common/env"
	"github.com/jackc/pgx/v4/pgxpool"
)

const DefaultMaxConnect = 1

type Config struct {
	ServiceName string
	Host        string
	User        string
	Name        string
	Password    string
	MaxConns    int
	Port        int
	SearchPath  string
}

// ConfigEnv returns config, all settings will be pulled from environment variables
func ConfigFromEnv(serviceName string) Config {
	return Config{
		ServiceName: serviceName,
		Host:        env.GetString("DB_HOST", "localhost"),
		User:        env.GetString("DB_USER", "postgres"),
		Name:        env.MustGetString("DB_NAME"),
		Password:    env.GetString("DB_PASSWORD", ""),
		Port:        env.GetInt("DB_PORT", 5432),
		SearchPath:  env.GetString("DB_SEARCH_PATH", ""),
		MaxConns:    DefaultMaxConnect,
	}
}

// connectionConfig returns the database connection config and error
func connectionConfig(cfg *Config) (*pgxpool.Config, error) {
	connStr := fmt.Sprintf("host=%s user=%s dbname=%s password=%s port=%d pool_max_conns=%d",
		cfg.Host,
		cfg.User,
		cfg.Name,
		cfg.Password,
		uint16(cfg.Port),
		cfg.MaxConns,
	)

	if len(cfg.SearchPath) > 0 || len(cfg.ServiceName) > 0 {
		if len(cfg.SearchPath) > 0 {
			connStr = fmt.Sprintf("%s search_path=%s", connStr, cfg.SearchPath)
		}
		if len(cfg.ServiceName) > 0 {
			connStr = fmt.Sprintf("%s application_name=%s", connStr, cfg.ServiceName)
		}
	}

	config, err := pgxpool.ParseConfig(connStr)
	config.ConnConfig.PreferSimpleProtocol = true
	return config, err
}

// Conn returns pgx connection pool and error
func Conn(ctx context.Context, cfg *Config) (*pgxpool.Pool, error) {
	config, err := connectionConfig(cfg)
	if err != nil {
		return nil, err
	}
	return pgxpool.ConnectConfig(ctx, config)
}
