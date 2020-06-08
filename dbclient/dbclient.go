package dbclient

import (
	"github.com/HomesNZ/env"
	"github.com/jackc/pgx"
)

// connectionConfig returns the database connection config
func connectionConfig(service string) pgx.ConnConfig {

	config := pgx.ConnConfig{
		Host:     env.GetString("DB_HOST", "localhost"),
		Database: env.MustGetString("DB_NAME"),
		User:     env.GetString("DB_USER", "postgres"),
		Password: env.GetString("DB_PASSWORD", ""),
		Port:     uint16(env.GetInt("DB_PORT", 5432)),
	}

	searchPath := env.GetString("DB_SEARCH_PATH", "")
	if len(searchPath) > 0 || len(service) > 0 {
		runtimeParams := map[string]string{}
		if len(searchPath) > 0 {
			runtimeParams["search_path"] = searchPath
		}
		if len(service) > 0 {
			runtimeParams["application_name"] = service
		}
		config.RuntimeParams = runtimeParams
	}
	return config
}

// Conn returns pgx connection pool and error
func Conn(service string, maxConnections int) (c *pgx.ConnPool, err error) {
	poolConfig := pgx.ConnPoolConfig{
		ConnConfig:     connectionConfig(service),
		MaxConnections: maxConnections,
	}
	return pgx.NewConnPool(poolConfig)
}
