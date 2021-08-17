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
			cfg := &Config{AccessKeyID: "key-id", SecretAccessKey: "secret", BucketName: "test-bucket"}
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("does not return an error for NewFromEnv, because it'll set default port and host", func() {
			accessKeyID := "key-id"
			secretAccessKey := "secret"
			bucketName := "test-bucket"
			os.Setenv("AWS_ACCESS_KEY_ID", accessKeyID)
			defer os.Unsetenv("AWS_ACCESS_KEY_ID")
			os.Setenv("AWS_SECRET_ACCESS_KEY", secretAccessKey)
			defer os.Unsetenv("AWS_SECRET_ACCESS_KEY")
			os.Setenv("AWS_S3_BUCKET", bucketName)
			defer os.Unsetenv("AWS_S3_BUCKET")
			cfg := NewFromEnv()
			err := cfg.Validate()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
