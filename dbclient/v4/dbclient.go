package dbclient

import (
	"context"
	"fmt"

	"github.com/HomesNZ/go-common/dbclient/v4/config"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

func New(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "DB")
	}
	config, err := connectionConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "DB")
	}
	return pgxpool.ConnectConfig(ctx, config)
}

func NewFromEnv(ctx context.Context) (*pgxpool.Pool, error) {
	cfg := config.NewFromEnv()
	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "DB")
	}
	config, err := connectionConfig(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "DB")
	}
	return pgxpool.ConnectConfig(ctx, config)
}

// connectionConfig returns the database connection config and error
func connectionConfig(cfg *config.Config) (*pgxpool.Config, error) {
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
	return config, errors.Wrap(err, "DB failed to parse config")
}
