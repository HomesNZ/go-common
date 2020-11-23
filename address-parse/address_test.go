package addressParse

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address", func() {
	newParser := func() (*Parser, error) {
		return New(
			StreetTypeAbbreviations(map[string]string{
				"RD": "ROAD",
				"ST": "STREET",
			}),
			StreetDirectionAbbreviations(map[string]string{
				"SOUTH": "SOUTH",
				"NORTH": "NORTH",
				"EAST":  "EAST",
				"E":     "EAST",
				"WEST":  "WEST",
			}),
		)
	}

	addressEntries := []TableEntry{
		Entry("an empty string returns an empty string", "", ""),
		Entry("shorthand state highway with space and shorthand direction", "SH2 E", "State Highway 2 East"),
	}
	addresses := []string{
		"1586 Motueka River West Bank Road, Motueka Valley",
		"Pohangina Valley East Road",
		"Motueka River West Bank Road",
		"1701 State Highway 2 East, Nukuhou",
		"134 Awakino Point East Road, Awakino Point",
		"34 Lake Road, St",
		"34 Lake Road, St Ar",
		"34 Lake Road, St Arnaud",
	}
	for _, v := range addresses {
		addressEntries = append(addressEntries, Entry(v, v, v))
	}

	DescribeTable(".Address",
		func(input string, expected string) {
			p, err := newParser()
			Expect(err).ToNot(HaveOccurred())
			parsed := p.Address(input)
			Expect(parsed).ToNot(BeNil())
			Expect(parsed.String()).To(Equal(expected))
		},
		addressEntries...,
	)

	Describe(".Address", func() {
		Context("with a street address", func() {
			Context("with a fully-formed identifier with street alpha", func() {
				It("correctly parses the address", func() {
					p, err := newParser()
					Expect(err).ToNot(HaveOccurred())
					parsed := p.Address("Flat 5, 58B Fictional Rd, Fake Suburb, Faketown")

					Expect(parsed).ToNot(BeNil())
					Expect(parsed.AddressIdentifier).ToNot(BeNil())
					Expect(parsed.AddressIdentifier.Len).To(Equal(12))
					Expect(parsed.Street).ToNot(BeNil())
					Expect(parsed.Street.Len).To(Equal(13))
					Expect(parsed.Street.Name).To(Equal("FICTIONAL"))
					Expect(parsed.Street.Type).To(Equal("ROAD"))
					Expect(parsed.Suburb).To(Equal("FAKE SUBURB"))
					Expect(parsed.City).To(Equal("FAKETOWN"))
				})
			})
			Context("with a fully-formed identifier with street alpha and postcode", func() {
				It("correctly parses the address", func() {
					p, err := newParser()
					Expect(err).ToNot(HaveOccurred())
					parsed := p.Address("Flat 5, 58B Fictional Rd, Fake Suburb, Faketown 1023")

					Expect(parsed).ToNot(BeNil())
					Expect(parsed.AddressIdentifier).ToNot(BeNil())
					Expect(parsed.AddressIdentifier.Len).To(Equal(12))
					Expect(parsed.Street).ToNot(BeNil())
					Expect(parsed.Street.Len).To(Equal(13))
					Expect(parsed.Street.Name).To(Equal("FICTIONAL"))
					Expect(parsed.Street.Type).To(Equal("ROAD"))
					Expect(parsed.Suburb).To(Equal("FAKE SUBURB"))
					Expect(parsed.City).To(Equal("FAKETOWN"))
					Expect(parsed.Postcode).To(Equal("1023"))
				})
			})
		})
		Context("with a street only", func() {
			It("correctly parses the address", func() {
				p, err := newParser()
				Expect(err).ToNot(HaveOccurred())
				parsed := p.Address("Kenya St")

				Expect(parsed).ToNot(BeNil())
				Expect(parsed.AddressIdentifier).To(BeNil())
				Expect(parsed.Street).ToNot(BeNil())
				Expect(parsed.Street.Len).To(Equal(8))
				Expect(parsed.Street.Name).To(Equal("KENYA"))
				Expect(parsed.Street.Type).To(Equal("STREET"))
				Expect(parsed.Suburb).To(Equal(""))
				Expect(parsed.City).To(Equal(""))
			})
		})
		Context("with a street only with no street type or direction", func() {
			It("correctly parses the address", func() {
				p, err := newParser()
				Expect(err).ToNot(HaveOccurred())
				parsed := p.Address("118 funnystreet")

				Expect(parsed).ToNot(BeNil())
				Expect(parsed.AddressIdentifier).ToNot(BeNil())
				Expect(parsed.AddressIdentifier.StreetNumber).To(Equal(int64(118)))
				Expect(parsed.Street).ToNot(BeNil())
				Expect(parsed.Street.Name).To(Equal("FUNNYSTREET"))
				Expect(parsed.Suburb).To(Equal(""))
				Expect(parsed.City).To(Equal(""))
			})
		})
	})
})
