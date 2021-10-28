package config

import (
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
			cfg := &Config{}
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
		})
		It("does not return an error", func() {
			cfg := &Config{ServiceName: "service", Name: "db"}
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return an error for NewFromEnv, because it'll set default port and host", func() {
			serviceName := "service"
			name := "db"
			os.Setenv("SERVICE_NAME", serviceName)
			defer os.Unsetenv("SERVICE_NAME")
			os.Setenv("DB_NAME", name)
			defer os.Unsetenv("DB_NAME")
			cfg := NewFromEnv()
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
