package redis

import (
	"github.com/gomodule/redigo/redis"
)

// Get returns a key value pair from redis.
func (c cache) Get(key string) (string, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))

	return reply, err
}

func (c cache) Exists(key string) (bool, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.Bool(conn.Do("EXISTS", key))

	return reply, err
}

