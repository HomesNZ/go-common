package redis

import (
	"time"
)

// Set adds a new key value pair to the redis cache.
func (c Cache) Set(key, val string) error {
	conn := c.Conn()
	defer conn.Close()

	_, err := conn.Do("SET", key, val)
	if err != nil {
		return err
	}

	return nil
}

// SetExpiry adds a new key value pair to the redis cache with expire time in seconds
func (c Cache) SetExpiry(key, val string, expireTime int) error {
	conn := c.Conn()
	defer conn.Close()

	_, err := conn.Do("SETEX", key, expireTime, val)
	if err != nil {
		return err
	}

	return nil
}
// SetExpiryTime adds a new key value pair to the redis cache with expire time in time.Time
func (c Cache) SetExpiryTime (key,val string, expireTime time.Time) error{
	conn := c.Conn()
	defer conn.Close()
	expire := expireTime.Unix() - time.Now().Unix()

	_, err := conn.Do("SETEX", key, int(expire), val)
	if err != nil {
		return err
	}

	return nil
}
