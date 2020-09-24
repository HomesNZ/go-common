package redis

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HomesNZ/go-common/env"
	log "github.com/sirupsen/logrus"
	"github.com/cenkalti/backoff"

	"github.com/gomodule/redigo/redis"
)

const (
	dbNumber = 0
)

var (
	// ConnBackoffTimeout is the duration before the backoff will timeout
	ConnBackoffTimeout = time.Duration(30) * time.Second

	// ErrUnableToConnectToRedis is raised when a connection to redis cannot be established.
	ErrUnableToConnectToRedis = errors.New("Unable to connect to redis")

	pool *redis.Pool

	once sync.Once
)

// Cache is a pool of connections to a redis cache
type Cache struct {
	Pool *redis.Pool
}

// Conn returns an active connection to the cache
func (c Cache) Conn() redis.Conn {
	return c.Pool.Get()
}

func addr() string {
	return fmt.Sprintf("%s:%s", env.MustGetString("REDIS_HOST"), env.MustGetString("REDIS_PORT"))
}

// CacheConn initializes (if not already initialized) and returns a connection to the redis cache
func CacheConn() Cache {
	once.Do(InitConnection)

	return Cache{
		Pool: pool,
	}
}

// SetConnection triggers the once lock, and returns a pool with the current connection
func SetConnection(c redis.Conn) Cache {
	once.Do(func() {})

	redisPool := &redis.Pool{
		//MaxIdle: 3,
		//IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return c, nil
		},
	}

	return Cache{
		Pool: redisPool,
	}
}

// InitConnection initializes a new redis cache connection pool.
func InitConnection() {
	//Creates a pool of connections to redis
	redisPool := &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		//Dial is an anonymous function which returns a redis.Conn
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr())
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	err := verifyConnection(redisPool.Get())
	if err != nil {
		log.WithError(err).Error(err)
	}

	pool = redisPool
}

// verifyConnection pings redis to verify a connection is established. If the connection cannot be established, it will
// retry with an exponential back off.
func verifyConnection(c redis.Conn) error {
	log.Infof("Attempting to connect to redis: %s", addr())

	pingDB := func() error {
		_, err := c.Do("PING")
		return err
	}

	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = ConnBackoffTimeout

	err := backoff.Retry(pingDB, expBackoff)
	if err != nil {
		log.Warning(err)
		log.Fatal(ErrUnableToConnectToRedis)
	}

	log.Info("Connected to redis")

	return nil
}
