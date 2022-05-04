package redis

import (
	"github.com/gomodule/redigo/redis"
)

// Delete removes a key from redis and returns its value
func (c cache) Delete(key string) (string, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	_, err = conn.Do("DEL", key)
	if err != nil {
		return "", err
	}

	return reply, err
}
