package redis

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	// ConnBackoffTimeout is the duration before the backoff will timeout
	ConnBackoffTimeout = time.Duration(30) * time.Second
)

type Cache interface {
	Delete(key string) (string, error)
	Get(key string) (string, error)
	Exists(key string) (bool, error)
	Set(key, val string) error
	SetExpiry(key, val string, expireTime int) error
	SetExpiryTime(key, val string, expireTime time.Time) error
	Subscribe(subscription string, handleResponse func(interface{}))
}

type cache struct {
	Pool *redisc.Cluster
	cfg  *Config
}

// Conn returns an active connection to the cache
func (c cache) Conn() redis.Conn {
	return c.Pool.Get()
}

func addr(cfg *Config) string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

func NewCache(cfg *Config) (*cache, error) {
	redisPool := &redisc.Cluster{
		CreatePool:   createPool,
		StartupNodes: []string{addr(cfg)},
	}
	log.Infof("Attempting to connect to redis: %s", addr(cfg))
	err := verifyConnection(redisPool.Get())
	if err != nil {
		return nil, errors.Wrapf(err, "Unable to connect to redis")
	}
	log.Info("Connected to redis")
	return &cache{Pool: redisPool, cfg: cfg}, nil
}

func createPool(address string, options ...redis.DialOption) (*redis.Pool, error) {
	return &redis.Pool{
		IdleTimeout: 60 * time.Second,
		// Dial is an anonymous function which returns a redis.Conn
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}, nil

}

// verifyConnection pings redis to verify a connection is established. If the connection cannot be established, it will
// retry with an exponential back off.
func verifyConnection(c redis.Conn) error {
	pingDB := func() error {
		_, err := c.Do("PING")
		return err
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = ConnBackoffTimeout

	return backoff.Retry(pingDB, expBackoff)
}
