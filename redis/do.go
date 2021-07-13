package redis

func (c cache) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := c.Conn()
	defer conn.Close()
	return conn.Do(commandName, args)
}
