package config

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

func TestGeocoder(t *testing.T) {
	ginkgo.RegisterFailHandler(Fail)
	gomega.RunSpecs(t, "Config")
}

var _ = Describe("Config", func() {
	Describe("#validate", func() {
		It("returns an error", func() {
			cfg := &Config{}
			err := cfg.Validate()
			Expect(err).To(HaveOccurred())
		})

		It("does not return an error for NewFromEnv, because it'll set default port and host", func() {
			apiKey := "apiKey"
			os.Setenv("BUGSNAG_API_KEY", apiKey)
			defer os.Unsetenv("BUGSNAG_API_KEY")
			cfg := NewFromEnv()
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
