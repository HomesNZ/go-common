package mailchimp_test

import (
	"github.com/HomesNZ/go-common/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCloudinary(t *testing.T) {
	config.InitLogger()

	RegisterFailHandler(Fail)
	RunSpecs(t, "MailChimp Suite")
}
