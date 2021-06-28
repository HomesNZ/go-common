package redis

import (
	"github.com/HomesNZ/go-common/redis/config"
	"github.com/mna/redisc"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func New(log *logrus.Entry, config config.Config) (Cache, error) {
	return newCache(log, config)
}

func NewFromEnv(log *logrus.Entry) (Cache, error) {
	cfg, err := config.NewFromEnv()
	if err != nil {
		return nil, err
	}

	return newCache(log, cfg)
}

func newCache(log *logrus.Entry, cfg config.Config) (Cache, error) {
	redisPool := &redisc.Cluster{
		CreatePool:   createPool,
		StartupNodes: []string{cfg.Addr()},
	}
	log.Infof("Attempting to connect to redis: %s", cfg.Addr())
	err := verifyConnection(redisPool.Get())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to redis")
	}
	log.Info("Connected to redis")
	return &cache{pool: redisPool, cfg: cfg}, nil
}
