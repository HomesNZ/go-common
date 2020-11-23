package redis

import (
	"os"

	"github.com/HomesNZ/go-common/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestRedis(t *testing.T) {
	config.InitLogger()

	RegisterFailHandler(Fail)
	RunSpecs(t, "Redis Suite")
}

var _ = BeforeSuite(func() {
	os.Setenv("REDIS_HOST", "redis.homes.dev")
	os.Setenv("REDIS_PORT", "6379")
})
