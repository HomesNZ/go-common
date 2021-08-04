package redis

import "github.com/gomodule/redigo/redis"

func (c cache) ListPush(listName string, val ...string) error {
	conn := c.Conn()
	defer conn.Close()
	values := make([]interface{}, 0, len(val)+1)
	values = append(values, listName)
	for _, v := range val {
		values = append(values, v)
	}
	_, err := conn.Do("LPUSH", values...)
	return err
}

func (c cache) ListLen(listName string) (int, error) {
	conn := c.Conn()
	defer conn.Close()
	return redis.Int(conn.Do("LLEN", listName))
}
func (c cache) ListPop(listName string) ([]string, error) {
	conn := c.Conn()
	defer conn.Close()
	return redis.Strings(conn.Do("LPOP", listName))
}

func (c cache) ListValues(listName string) ([]string, error) {
	conn := c.Conn()
	defer conn.Close()
	return redis.Strings(conn.Do("LRANGE", listName, 0, -1))
}
