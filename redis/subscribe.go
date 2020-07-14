package redis

import (
	"github.com/sirupsen/logrus"
)

// Subscribe creates a subscription to the redis publishing system
//
// Messages will be passed into handleResponse.
// Subscribe will block forever, so a goroutine is recommended
func (c Cache) Subscribe(subscription string, handleResponse func(interface{})) {
	conn := c.Conn()
	defer conn.Close()

	_, err := conn.Do("PSUBSCRIBE", subscription)
	if err != nil {
		logrus.WithError(err).Error(err)
	}

	for err == nil {
		reply, err := conn.Receive()
		if err != nil {
			logrus.WithError(err).Fatal("Could not connect to redis ", err.Error())
		}

		handleResponse(reply)
	}
}
