package config

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGeocoder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config")
}

var _ = Describe("Config", func() {
	Describe("#validate", func() {
		It("returns an error", func() {
			err := config{}.validate()
			Expect(err).To(HaveOccurred())
		})
		It("does not return an error", func() {
			err := config{port: "1000", host: "host"}.validate()
			Expect(err).NotTo(HaveOccurred())
		})
		It("does not return an error for NewFromEnv, because it'll set default port and host", func() {
			cfg, err := NewFromEnv()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Addr()).To(Equal("localhost:6379"))
		})

		It("does not return an error for NewFromEnv, because it'll set default port and host", func() {
			host := "www.homes.co.nz"
			port := "5000"
			os.Setenv("REDIS_HOST", host)
			defer os.Unsetenv("REDIS_HOST")
			os.Setenv("REDIS_PORT", port)
			defer os.Unsetenv("REDIS_PORT")
			cfg, err := NewFromEnv()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Addr()).To(Equal(fmt.Sprintf("%s:%s", host, port)))
		})
	})
})
