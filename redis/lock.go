package redis

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

// Lockable interface defines the
// way we can compare and decide if redis lock exist
// and whether existing lock
// need to be updated
type Lockable interface {
	Updated() string
	Key() string
}

// IsProcessed checks if the lockable entity
// were processed by comparing it with
// corresponding redis key and updated field
func (c cache) IsProcessed(lockable Lockable) (bool, error) {
	resp, err := c.Get(lockable.Key())
	// ErrNil if the key is not found
	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, errors.Wrap(err, "error getting key")
	}
	return lockable.Updated() == resp, nil
}

func (c cache) MarkProcessed(lockable Lockable) error {
	err := c.Set(lockable.Key(), lockable.Updated())
	return err
}
