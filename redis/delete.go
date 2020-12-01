package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

// Delete removes a key from redis and returns its value
func (c Cache) Delete(key string) (string, error) {
	conn := c.Conn()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		logrus.WithError(err).Error(err)
		return "", err
	}

	_, err = conn.Do("DEL", key)
	if err != nil {
		logrus.WithError(err).Error(err)
		return "", err
	}

	return reply, err
}
