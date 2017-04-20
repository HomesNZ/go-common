package mailchimp_test

import (
	. "github.com/HomesNZ/go-common/mailchimp"
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MailChimp.Config", func() {
	DescribeTable(".Validate",
		func(config Config, expected error) {
			if expected == nil {
				Expect(config.Validate()).To(BeNil())
			} else {
				Expect(config.Validate()).To(Equal(expected))
			}
		},
		Entry("returns an error if data center is empty", Config{APIKey: "abc123"}, ErrConfigEmptyDataCenter),
		Entry("returns an error if api key is empty", Config{DataCenter: "abc123"}, ErrConfigEmptyAPIKey),
		Entry("returns an error if data center is empty", Config{DataCenter: "abc123", APIKey: "abc123"}, nil),
	)
})
