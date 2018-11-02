package redis_test

import (
	. "github.com/HomesNZ/go-common/redis"
	"github.com/rafaeljusto/redigomock"

	"time"

	"github.com/gomodule/redigo/redis"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redis", func() {
	Describe(".Conn", func() {
		Context("with an external connection [integration]", func() {
			It("returns an active connection to Redis", func() {
				c := CacheConn()

				_, err := c.Conn().Do("PING")

				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("with an mocked connection", func() {
			It("returns an active connection to Redis", func() {
				// TODO (nicday): Commented out to get tests working for travis. This test expects an actual connection to reset
				//  to after the test has finished. In reality we don't need to reset the connection here, we just don't want to
				// create interacting tests.
				//
				// existingRedis := CacheConn()
				// // defer statments are LIFO, so this is run after the mock connection is closed.
				// defer SetConnection(existingRedis.Conn)

				mock := redigomock.NewConn()
				mock.Command("PING")

				c := SetConnection(mock)
				_, err := c.Conn().Do("PING")

				mock.Close()
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe(".Subscribe", func() {
		It("Outputs the expected key", func() {
			c := CacheConn()

			go c.Subscribe("*:test", func(value interface{}) {
				Expect(value).To(Equal("name:test"))
			})

			c.SetExpiry("name:test", "value", 1)
			time.Sleep(2)

		})
	})
})

type Pool interface {
	Get() redis.Conn
}
