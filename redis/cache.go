package redis

import (
	"github.com/HomesNZ/go-common/redis/config"
	"github.com/cenkalti/backoff"
	"github.com/gomodule/redigo/redis"
	"github.com/mna/redisc"
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
	pool *redisc.Cluster
	cfg  config.Config
}

// Conn returns an active connection to the cache
func (c cache) Conn() redis.Conn {
	return c.pool.Get()
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