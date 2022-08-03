package redis

import (
	"github.com/gomodule/redigo/redis"
)

// For compatibility
func (c cache) Get(key string) (string, error) {
	return c.GetString(key)
}

func (c cache) GetString(key string) (string, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))

	return reply, err
}

func (c cache) GetBool(key string) (bool, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.Bool(conn.Do("GET", key))

	return reply, err
}

func (c cache) Exists(key string) (bool, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.Bool(conn.Do("EXISTS", key))

	return reply, err
}
