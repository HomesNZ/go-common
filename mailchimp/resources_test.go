package mailchimp

import (
	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MailChimp.Resources", func() {
	DescribeTable("#formatUrl",
		func(i testInput, e testOutput) {
			o := testOutput{}
			o.str, o.err = formatUrl(i.url, i.params)
			if e.err == nil {
				Expect(o.err).To(BeNil())
			} else {
				Expect(o.err).To(HaveOccurred())
			}
			Expect(o.str).To(Equal(e.str))
		},
		Entry(
			"returns an error if the resource does not exist",
			testInput{"some-random-string", nil},
			testOutput{"", ErrUnknownResource},
		),
		/*
			    Entry(
					  "returns an error if the template could not be executed",
					   input{"", nil},
					    output{"",ErrBadlyFormattedURL}
					  ),
		*/
		Entry(
			"returns a url without parameters as-is",
			testInput{PostAuthorizedApps, nil},
			testOutput{PostAuthorizedApps, nil},
		),
		Entry(
			"returns a formatted url with a single parameter",
			testInput{GetAuthorizedAppsAppID, AuthorizedAppParams{AppID: "12345"}},
			testOutput{"/authorized-apps/12345", nil},
		),
		Entry(
			"returns a formatted url with multiple parameters",
			testInput{
				GetAutomationsWorkflowIDEmailsWorkflowEmailIDQueueSubscriberHash,
				AutomationParams{
					WorkflowID:      "param1",
					WorkflowEmailID: "param2",
					SubscriberHash:  "param3",
				},
			},
			testOutput{"/automations/param1/emails/param2/queue/param3", nil},
		),
	)
})

type testInput struct {
	url    string
	params interface{}
}

type testOutput struct {
	str string
	err error
}
