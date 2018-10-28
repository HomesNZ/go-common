package nzAddress_test

import (
	. "github.com/HomesNZ/go-common/nz_address"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Address", func() {
	DescribeTable(".Display", func(input Address, expected string) {
		Expect(input.Display()).To(Equal(expected))
	},
		Entry(
			"a full set of address identifiers",
			Address{
				UnitIdentifier:   "123",
				StreetNumber:     5,
				StreetAlpha:      "A",
				StreetNumberHigh: 10,
			},
			"123/5A-10",
		),
		Entry(
			"street number range",
			Address{
				StreetNumber:     5,
				StreetNumberHigh: 10,
			},
			"5-10",
		),
		Entry(
			"street number range with unit identifier",
			Address{
				UnitIdentifier:   "123",
				StreetNumber:     5,
				StreetNumberHigh: 10,
			},
			"123/5-10",
		),
		Entry(
			"a full set of address identifiers and street name/type",
			Address{
				UnitType:       "FLAT",
				UnitIdentifier: "123",
				StreetNumber:   5,
				StreetAlpha:    "A",
				StreetName:     "Cambridge",
				StreetType:     "Terrace",
			},
			"Flat 123/5A Cambridge Terrace",
		),
		Entry(
			"a full set of address identifiers passed as 'identifier' and street name/type",
			Address{
				Identifier: "FLAT 123/5A",
				StreetName: "Cambridge",
				StreetType: "Terrace",
			},
			"FLAT 123/5A Cambridge Terrace",
		),
		Entry(
			"a full set of address identifiers passed as 'identifier' AND separate address fields and street name/type",
			Address{
				Identifier:     "FLAT 123/5A",
				UnitType:       "FLAT",
				UnitIdentifier: "123",
				StreetNumber:   5,
				StreetAlpha:    "A",
				StreetName:     "Cambridge",
				StreetType:     "Terrace",
			},
			"FLAT 123/5A Cambridge Terrace",
		),
		Entry(
			"building name",
			Address{
				BuildingName: "Homes",
				StreetNumber: 123,
				StreetName:   "Cambridge",
				StreetType:   "Terrace",
			},
			"Homes, 123 Cambridge Terrace",
		),
		Entry(
			"RD number",
			Address{
				StreetNumber: 123,
				StreetName:   "Cambridge",
				StreetType:   "Terrace",
				RDNumber:     "3c",
				City:         "Wellington",
			},
			"123 Cambridge Terrace, RD 3C, Wellington",
		),
		Entry(
			"city and suburb",
			Address{
				Suburb: "Brooklyn",
				City:   "Wellington",
			},
			"Brooklyn, Wellington",
		),
		Entry(
			"duplicate city and suburb",
			Address{
				Suburb: "Wellington",
				City:   "Wellington",
			},
			"Wellington",
		),
		Entry(
			"street name and direction without type",
			Address{
				StreetName:      "State Highway 2",
				StreetDirection: "East",
			},
			"State Highway 2 East",
		),
		Entry(
			"full address",
			Address{
				BuildingName:     "Homes House",
				UnitType:         "Unit",
				UnitIdentifier:   "5",
				StreetNumber:     123,
				StreetNumberHigh: 134,
				StreetAlpha:      "B",
				StreetName:       "CAMBRIDGE",
				StreetType:       "TERRACE",
				RDNumber:         "3a",
				Suburb:           "BROOKLYN",
				City:             "WELLINGTON",
				Postcode:         1234,
			},
			"Unit 5 Homes House, 123B-134 Cambridge Terrace, Brooklyn, RD 3A, Wellington",
		),
		Entry(
			"full address with duplicate suburb/city",
			Address{
				BuildingName:     "Homes House",
				UnitType:         "Unit",
				UnitIdentifier:   "5",
				StreetNumber:     123,
				StreetNumberHigh: 134,
				StreetAlpha:      "B",
				StreetName:       "CAMBRIDGE",
				StreetType:       "TERRACE",
				RDNumber:         "3a",
				Suburb:           "WELLINGTON",
				City:             "WELLINGTON",
				Postcode:         1234,
			},
			"Unit 5 Homes House, 123B-134 Cambridge Terrace, RD 3A, Wellington",
		),
	)
	Describe("DisplayWithPostcode", func() {
		It("appends postcode to the end of the display address", func() {
			input := Address{
				BuildingName:   "Homes House",
				UnitType:       "UNIT",
				UnitIdentifier: "5",
				StreetNumber:   123,
				StreetAlpha:    "B",
				StreetName:     "CAMBRIDGE",
				StreetType:     "TERRACE",
				RDNumber:       "3a",
				Suburb:         "BROOKLYN",
				City:           "WELLINGTON",
				Postcode:       1234,
			}.DisplayWithPostcode()

			expected := "Unit 5 Homes House, 123B Cambridge Terrace, Brooklyn, RD 3A, Wellington 1234"

			Expect(input).To(Equal(expected))
		})
	})
})
