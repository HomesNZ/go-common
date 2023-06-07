package redis

import (
	"github.com/HomesNZ/go-common/redis/config"
	"github.com/mna/redisc"
	"github.com/pkg/errors"
)

func New(config config.Config) (Cache, error) {
	return newCache(config)
}

func NewFromEnv() (Cache, error) {
	cfg, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newCache(cfg)
}

func newCache(cfg config.Config) (Cache, error) {
	redisPool := &redisc.Cluster{
		CreatePool:   createPool,
		StartupNodes: []string{cfg.Addr()},
	}
	err := verifyConnection(redisPool.Get())
	if err != nil {
		return nil, errors.WithMessage(err, "Unable to connect to redis")
	}

	return &cache{pool: redisPool, cfg: cfg}, nil
}
