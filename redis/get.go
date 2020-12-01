package redis

import (
	"github.com/gomodule/redigo/redis"
)

// Get returns a key value pair from redis.
func (c Cache) Get(key string) (string, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))

	return reply, err
}
