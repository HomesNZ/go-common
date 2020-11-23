package addressParse

import (
	"strconv"
	"strings"

	"github.com/HomesNZ/go-common/nz_address"
)

// ParsedAddress represents a parsed address in its entirety
type ParsedAddress struct {
	AddressIdentifier *ParsedAddressIdentifier
	Street            *ParsedStreet
	Suburb            string
	City              string
	Postcode          string
	Raw               string
}

// Address parses an entire address string and returns the separate parts
/*
	// TODO: finish fully parsing address identifiers
	input: "Flat 5, 58B Fictional Rd, Fake Suburb, Faketown"
	output: *ParsedAddress{
		AddressIdentifier *ParsedAddressIdentifier{
			Type:					  "Flat",
			UnitIdentifier: "5",
			StreetNumber:   "58",
			StreetAlpha:    "B",
		},
		Street            *ParsedStreet{
			Name: "Fictional",
			Type: "Rd",
		},
		Suburb            "Fake Suburb",
		City              "Faketown",
	}

	input: "Kenya St"
	output: *ParsedAddress{
		AddressIdentifier *ParsedAddressIdentifier
		Street            *ParsedStreet{
			Name: "Kenya",
			Type: "Street",
		}
		Suburb            "",
		City              "",
	}

	input: "Fake Suburb"
	output: *ParsedAddress{
		AddressIdentifier nil,
		Street            nil,
		Suburb            "Fake Suburb",
		City              "",
	}
*/
func (p Parser) Address(address string) *ParsedAddress {
	parsed := ParsedAddress{Raw: address}

	// Extract the address identifier
	parsed.AddressIdentifier = p.AddressIdentifier(address)
	if parsed.AddressIdentifier != nil {
		// Check if the whole string was considered an address identifier. Chances
		// are it's actually a street name, because people don't search solely for
		// an address identifier. AddressIdentifier method does not have this
		// context so is unable to make that call.
		if parsed.AddressIdentifier.Len >= len(strings.Trim(address, " ")) {
			parsed.AddressIdentifier = nil
		} else {
			start := parsed.AddressIdentifier.Start + parsed.AddressIdentifier.Len
			address = address[start:]
		}
	}

	// Parse the street name
	parsed.Street = p.Street(address)
	if parsed.Street != nil {
		start := parsed.Street.Start + parsed.Street.Len
		address = address[start:]
	}

	// Match city/suburb
	parsed.Suburb, parsed.City, parsed.Postcode = p.SuburbAndCity(address)

	return &parsed
}

// NZAddress converts the ParsedAddress to a nzAddress.Address and returns it.
func (a ParsedAddress) NZAddress() nzAddress.Address {
	postcode, _ := strconv.Atoi(a.Postcode)
	addr := nzAddress.Address{
		Suburb:   a.Suburb,
		City:     a.City,
		Postcode: postcode,
	}
	if a.AddressIdentifier != nil {
		addr.UnitIdentifier = a.AddressIdentifier.UnitIdentifier
		addr.StreetNumber = int(a.AddressIdentifier.StreetNumber)
		addr.StreetAlpha = a.AddressIdentifier.StreetAlpha
	}
	if a.Street != nil {
		addr.StreetName = a.Street.Name
		addr.StreetType = a.Street.Type
		addr.StreetDirection = a.Street.Direction
	}
	return addr
}

func (a ParsedAddress) String() string {
	return a.NZAddress().Display()
}

// Title returns the "identifying" part of the address to match a
// typeahead_search record. If it's parsed as a street address, then that's
// just the address identifier and street name/type/direction.
func (a ParsedAddress) Title() string {
	nza := nzAddress.Address{}
	if a.AddressIdentifier != nil {
		nza.UnitIdentifier = a.AddressIdentifier.UnitIdentifier
		nza.StreetNumber = int(a.AddressIdentifier.StreetNumber)
		nza.StreetAlpha = a.AddressIdentifier.StreetAlpha
		if a.Street != nil {
			nza.StreetName = a.Street.Name
			nza.StreetType = a.Street.Type
			nza.StreetDirection = a.Street.Direction
			return nza.Display()
		}
	}
	return a.Raw
}
