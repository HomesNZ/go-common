package addressParse

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("AddressIdentifier", func() {
	DescribeTable(".AddressIdentifier",
		func(address string, expected string) {
			p, err := New()
			Expect(err).ToNot(HaveOccurred())
			res := p.AddressIdentifier(address)
			fmt.Println(address, "||||", res, "||||", expected)
			if expected == "" {
				Expect(res).To(BeNil())
			} else {
				Expect(res).ToNot(BeNil())
				Expect(res.String()).To(Equal(expected))
			}
		},
		Entry(
			"state highway address",
			"1234 State Highway 12, Some Suburb",
			"1234",
		),
		Entry(
			"state highway with suburb",
			"State Highway 12, Some Suburb",
			"",
		),
		Entry(
			"state highway without suburb",
			"State Highway 46",
			"",
		),
		Entry(
			"short form state highway without spaces",
			"SH46",
			"",
		),
		Entry(
			"short form state highway with spaces",
			"S H 46",
			"",
		),
		Entry(
			"state highway address with unit identifier",
			"123/1234 State Highway 12, Some Suburb",
			"123/1234",
		),
		Entry(
			"unit identifier and street number range",
			"C4 23-25",
			"C4/23",
		),
		Entry(
			"address with numeric unit identifier and street number and alpha",
			"1/179A Birkdale Road,  Birkdale, Auckland 0626",
			"1/179A",
		),
		Entry(
			"address with full-form unit type and numeric unit identifier and street number and alpha",
			"Unit 53, 18A Cuba Street, Te Aro, Wellington",
			"Unit 53/18A",
		),
		Entry(
			"address with full-form unit type and numeric unit identifier and street number",
			"Unit 53, 18 Cuba Street, Te Aro, Wellington",
			"Unit 53/18",
		),
		Entry(
			"address with full-form unit type and numeric unit identifier and street number",
			"22 Grampian Road, St Heliers, Auckland 1041",
			"22",
		),
		Entry(
			"address with street number",
			"18 Cuba Street, Te Aro, Wellington",
			"18",
		),
		Entry(
			"address with street number and uppercase alpha",
			"18A Cuba Street, Te Aro, Wellington",
			"18A",
		),
		Entry(
			"address with street number and lowercase alpha",
			"18a Cuba Street, Te Aro, Wellington",
			"18A",
		),
		Entry(
			"address with numeric unit identifier with street number",
			"23/18 Cuba Street, Te Aro, Wellington",
			"23/18",
		),
		Entry(
			"address with single-digit street number",
			"1 Cuba Street, Te Aro, Wellington",
			"1",
		),
		Entry(
			"address with alphanumeric unit identifier and with street number and alpha",
			"23A/18F Cuba Street, Te Aro, Wellington",
			"23A/18F",
		),
		Entry(
			"street and city",
			"Cuba Street, Wellington",
			"",
		),
		Entry(
			"city only",
			"Wellington",
			"",
		),
		Entry(
			"street only",
			"Cuba Street",
			"",
		),
		Entry(
			"suburb and city",
			"Te Aro, Wellington",
			"",
		),
		Entry(
			"suburb only",
			"Te Aro",
			"",
		),
		// Address identifiers only
		Entry(
			"concatenates a char sub number",
			"123A",
			"123A",
		),
		Entry(
			"concatenates a numeric sub number with a separator",
			"2/123",
			"2/123",
		),
		Entry(
			"uses the standard separator if there is a non-standard one present in the sub number",
			"1\\123",
			"1/123",
		),
		Entry(
			"does not concatenate if sub number is empty",
			"123",
			"123",
		),
		Entry(
			"is an empty string if both number and sub number are empty",
			"",
			"",
		),
		Entry(
			"correctly handles an alpha followed by a space and numeric",
			"B 23/3",
			"23/3B",
		),
		Entry(
			"correctly handles a flat number",
			"F456/123",
			"Flat 456/123",
		),
		Entry(
			"correctly handles a unit number",
			"U456/123",
			"Unit 456/123",
		),
		Entry(
			"correctly handles unit id with a numeric followed by an alpha",
			"12A/123",
			"12A/123",
		),
		Entry(
			"accepts any unit identifier starting with a numeric character",
			"1 ABC DEF G-sdssds/123",
			"1 ABC DEF G-SDSSDS/123",
		),
		Entry(
			"correctly handles unit id with an alpha followed by a numeric",
			"A1/123",
			"A1/123",
		),
		Entry(
			"abbreviates unit identifier and prefixes the street number with the abbreviated unit identifier.",
			"Unit 1/123",
			"Unit 1/123",
		),
		Entry(
			"abbreviates unit identifier with no spaces and prefixes the street number with the abbreviated unit identifier.",
			"Unit1/123",
			"Unit 1/123",
		),
		Entry(
			"abbreviates unit identifier and prefixes the street number with the abbreviated unit identifier.",
			"Flat 1/123",
			"Flat 1/123",
		),
		Entry(
			"abbreviates unit identifier with no spaces and prefixes the street number with the abbreviated unit identifier.",
			"Flat 1/123",
			"Flat 1/123",
		),
		Entry(
			"villa",
			"Villa 1/123",
			"Villa 1/123",
		),
		Entry(
			"villas",
			"Villas 1-2/123",
			"Villa 1-2/123",
		),
		Entry(
			"lot number",
			"Lot 1, 23",
			"Lot 1/23",
		),
		Entry(
			"lot number, no street number",
			"Lot 1",
			"",
		),
		Entry(
			"lot(s) number",
			"Lots 1-2, 23",
			"Lot 1-2/23",
		),
		Entry(
			"space between street number and alpha",
			"3 A",
			"3",
		),
	)
})
