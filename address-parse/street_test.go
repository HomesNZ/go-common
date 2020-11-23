package addressParse_test

import (
	addressParse "github.com/HomesNZ/go-common/address-parse"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var Parser *addressParse.Parser

var _ = Describe("Street", func() {

	BeforeEach(func() {
		Parser, _ = addressParse.New(addressParse.StreetTypeAbbreviations(map[string]string{
			"AVENUE":    "AVENUE",
			"AVE":       "AVENUE",
			"ROAD":      "ROAD",
			"STREET":    "STREET",
			"ST":        "STREET",
			"HWY":       "HIGHWAY",
			"HIGHWAY":   "HIGHWAY",
			"RIVER":     "RIVER",
			"PROMENADE": "PROMENADE",
		}), addressParse.StreetDirectionAbbreviations(map[string]string{
			"SOUTH": "SOUTH",
			"NORTH": "NORTH",
			"EAST":  "EAST",
			"E":     "EAST",
			"WEST":  "WEST",
		}))
		//InitOnce.Do(func() {})

	})

	DescribeTable(".Street",
		func(input string, expected *addressParse.ParsedStreet) {
			parsed := Parser.Street(input)
			if expected == nil {
				Expect(parsed).To(BeNil())
			} else {
				Expect(parsed).ToNot(BeNil())
				Expect(parsed.Name).To(Equal(expected.Name))
				Expect(parsed.Type).To(Equal(expected.Type))
				Expect(parsed.Direction).To(Equal(expected.Direction))
				Expect(parsed.Len).To(Equal(expected.Len))
			}
		},
		Entry("nil", "", nil),
		Entry("type and direction before correct type", "MOTUEKA RIVER WEST BANK ROAD", &addressParse.ParsedStreet{Name: "MOTUEKA RIVER WEST BANK", Type: "ROAD", Len: 28}),
		Entry("state highway with direction", "STATE HIGHWAY 2 EAST", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2", Direction: "EAST", Len: 20}),
		Entry("lowercase state highway with direction", "State Highway 2 East", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2", Direction: "EAST", Len: 20}),
		Entry("state highway with alpha", "State Highway 2A", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2A", Len: 16}),
		Entry("state highway with alpha and direction", "State Highway 2A East", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2A", Direction: "EAST", Len: 21}),
		Entry("direction before type", "Owen Valley East Road", &addressParse.ParsedStreet{Name: "OWEN VALLEY EAST", Type: "ROAD", Len: 21}),
		Entry("shorthand state highway with space", "SH 2", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2", Len: 4}),
		Entry("shorthand state highway with space and shorthand direction", "SH2 E", &addressParse.ParsedStreet{Name: "STATE HIGHWAY 2", Direction: "EAST", Len: 5}),
	)

	Describe(".Street", func() {
		Context("with a normal street name", func() {
			It("correctly parses the street name without error", func() {
				parsedStreet := Parser.Street("Alexander Avenue")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ALEXANDER"))
				Expect(parsedStreet.Type).To(Equal("AVENUE"))
			})
		})
		Context("with a street name and no type", func() {
			It("correctly parses the street name without error", func() {
				parsedStreet := Parser.Street("Alexander")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ALEXANDER"))
				Expect(parsedStreet.Type).To(Equal(""))
				Expect(parsedStreet.Start).To(Equal(0))
				Expect(parsedStreet.Len).To(Equal(9))
			})
		})
		Context("using street name with an abbreviated type", func() {
			It("correctly parses the street name without error", func() {
				parsedStreet := Parser.Street("Alexander Ave")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ALEXANDER"))
				Expect(parsedStreet.Type).To(Equal("AVENUE"))
			})
		})
		Context("using street name with an multiple types detected", func() {
			It("only uses the last type found", func() {
				parsedStreet := Parser.Street("St Heliers Bay Road")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ST HELIERS BAY"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
			})
		})
		Context("using street name with an multiple directions and types detected", func() {
			It("only uses the last direction found", func() {
				parsedStreet := Parser.Street("Big Road South Small Road North")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("BIG ROAD SOUTH SMALL"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
				Expect(parsedStreet.Direction).To(Equal("NORTH"))
			})
		})
		Context("using a street name with suburb and city after it", func() {
			It("drops anything after the street type is found", func() {
				parsedStreet := Parser.Street("Kenya Street, Ngaio, Wellington")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("KENYA"))
				Expect(parsedStreet.Type).To(Equal("STREET"))
			})
		})
		Context("using a street name with superfluous spaces", func() {
			It("ignores the extra spaces", func() {
				parsedStreet := Parser.Street("  Kenya   Street  ")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("KENYA"))
				Expect(parsedStreet.Type).To(Equal("STREET"))
			})
		})
		Context("using a street name with a direction", func() {
			It("correctly parses the street name", func() {
				parsedStreet := Parser.Street("Kenya Street South")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("KENYA"))
				Expect(parsedStreet.Type).To(Equal("STREET"))
				Expect(parsedStreet.Direction).To(Equal("SOUTH"))
			})
		})
		Context("using a street name with suburb and no comma", func() {
			It("ignores any extra characters after the street type", func() {
				parsedStreet := Parser.Street("Kenya Street Ngaio")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("KENYA"))
				Expect(parsedStreet.Type).To(Equal("STREET"))
				Expect(parsedStreet.Direction).To(Equal(""))
			})
		})
		Context("using a full form state highway", func() {
			It("correctly parses the street name and type", func() {
				parsedStreet := Parser.Street("State Highway 56")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("STATE HIGHWAY 56"))
				Expect(parsedStreet.Type).To(Equal(""))
				Expect(parsedStreet.Direction).To(Equal(""))
			})
		})
		Context("using an abbreviated state highway", func() {
			It("correctly parses the street name and type", func() {
				parsedStreet := Parser.Street("State Hwy 56")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("STATE HIGHWAY 56"))
				Expect(parsedStreet.Type).To(Equal(""))
				Expect(parsedStreet.Direction).To(Equal(""))
			})
		})
		Context("using street with ST instead of SAINT", func() {
			It("doesn't swap ST for SAINT", func() {
				parsedStreet := Parser.Street("St Heliers Bay Road")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ST HELIERS BAY"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
			})
		})
		Context("using partial street with ST instead of SAINT", func() {
			It("doesn't swap ST for SAINT", func() {
				parsedStreet := Parser.Street("St Heli")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ST HELI"))
				Expect(parsedStreet.Type).To(Equal(""))
			})
		})
		Context("using SAINT with street with ST instead of SAINT", func() {
			It("swaps SAINT for ST", func() {
				parsedStreet := Parser.Street("Saint Heliers Bay Road")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ST HELIERS BAY"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
			})
		})
		Context("street name with suburb name", func() {
			It("strips out the separating comma", func() {
				parsedStreet := Parser.Street("Aranui Road, Mt Wellington")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("ARANUI"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
			})
		})
		Context("street name with a 'THE'", func() {
			It("parses as the street name with no street type", func() {
				parsedStreet := Parser.Street("The Promenade")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("THE PROMENADE"))
				Expect(parsedStreet.Type).To(Equal(""))
			})
		})
		Context("street name with direction before type", func() {
			It("adds the direction to the street name", func() {
				parsedStreet := Parser.Street("Owen Valley East Road")

				Expect(parsedStreet).ToNot(BeNil())
				Expect(parsedStreet.Name).To(Equal("OWEN VALLEY EAST"))
				Expect(parsedStreet.Type).To(Equal("ROAD"))
			})
		})
	})
	Describe("ParsedStreet.String", func() {
		Context("with a street name and type", func() {
			It("returns the correct string", func() {
				Expect(addressParse.ParsedStreet{
					Name:      "KENYA",
					Type:      "STREET",
					Direction: "",
				}.String()).To(Equal("Kenya Street"))
			})
		})
		Context("with a street name, type, and direction", func() {
			It("returns the correct string", func() {
				Expect(addressParse.ParsedStreet{
					Name:      "KENYA",
					Type:      "STREET",
					Direction: "WEST",
				}.String()).To(Equal("Kenya Street West"))
			})
		})
	})
})
