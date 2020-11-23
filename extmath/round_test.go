package extmath

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type roundInput struct {
	n  float64
	sf int
}

var _ = Describe("extmath", func() {
	DescribeTable(".Round",
		func(input roundInput, expected float64) {
			res := Round(input.n, input.sf)
			Expect(res).To(Equal(expected))
		},
		Entry("rounded number", roundInput{10.0, 5}, 10.0),
		Entry("round down to 2 sf", roundInput{10.12345, 2}, 10.12),
		Entry("round up to 2 sf", roundInput{10.98765, 2}, 10.99),
		Entry("round 2sf num to 2 sf", roundInput{10.98, 2}, 10.98),
		Entry("round down to 0 sf", roundInput{10.12, 0}, 10.0),
		Entry("round up to 0 sf", roundInput{10.98, 0}, 11.0),
		Entry("round negative up to 2 sf", roundInput{-10.12345, 2}, -10.12),
		Entry("round negative down to 2 sf", roundInput{-10.98765, 2}, -10.99),
		Entry("round negative up to 0 sf", roundInput{-10.12, 0}, -10.0),
		Entry("round negative down to 0 sf", roundInput{-10.98, 0}, -11.0),
	)
})
