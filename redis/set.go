package redis

import (
	"time"
)

// Set adds a new key value pair to the redis cache.
func (c cache) Set(key, val string) error {
	conn := c.Conn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		return err
	}

	return nil
}

// SetExpiry adds a new key value pair to the redis cache with expire time in seconds
func (c cache) SetExpiry(key, val string, expireTime int) error {
	conn := c.Conn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, expireTime, val)
	if err != nil {
		return err
	}

	return nil
}

// SetExpiryTime adds a new key value pair to the redis cache with expire time in time.Time
func (c cache) SetExpiryTime(key, val string, expireTime time.Time) error {
	conn := c.Conn()
	defer conn.Close()
	// convert the expiry time into the duration until expiry to conform to redis expectations
	expire := expireTime.Unix() - time.Now().Unix()

	_, err := conn.Do("SETEX", key, int(expire), val)
	if err != nil {
		return err
	}

	return nil
}

