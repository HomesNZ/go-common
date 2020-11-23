package addressParse

import (
	"github.com/pkg/errors"
	"strings"

	"github.com/HomesNZ/go-common/nz_address"
)

const (
	UnitTypeFlat  = "FLAT"
	UnitTypeLot   = "LOT"
	UnitTypeUnit  = "UNIT"
	UnitTypeVilla = "VILLA"
)

// ParsedAddressIdentifier is an address identifier parsed from an address string
type ParsedAddressIdentifier struct {
	StreetNumber     int64
	StreetNumberHigh int64
	StreetAlpha      string
	UnitType         string
	UnitIdentifier   string
	Start            int
	Len              int
}

// ParsedAddressIdentifiers is a slice of ParsedAddressIdentifier pointers
type ParsedAddressIdentifiers []*ParsedAddressIdentifier

// AddressIdentifier parses the address identifier fields from an address
// string. "address" can include street/suburb/city etc, or just the identifiers
// by themselves.
func (p Parser) AddressIdentifier(address string) *ParsedAddressIdentifier {
	if address == "" {
		return nil
	}

	// Convert address to uppercase, use that for matching to simplify everything.
	addressUpper := strings.ToUpper(address)

	// Check if any of the invalid matchers match. If any do, then this string
	// is certain to not contain any address identifiers.
	for _, m := range p.invalidMatchers() {
		if m.Invalid(addressUpper) {
			return nil
		}
	}

	// Run through all the matchers, and add any matches to the identifiers slice.
	identifiers := ParsedAddressIdentifiers{}
	for _, m := range p.matchers() {
		p.logger.Debugf("\nmatching '%s'...", m.MatcherName())
		identifier := m.Parse(addressUpper)
		if identifier == nil {
			continue
		}
		p.logger.Debugf("matched '%s': %#v", m.MatcherName(), *identifier)
		identifiers = append(identifiers, identifier)
	}

	// If none of the matchers were successful, stop here.
	if len(identifiers) == 0 {
		return nil
	}

	// Merge all the identifiers into one, remove any conflicting matches.
	final := p.mergeIdentifiers(&identifiers)

	// If no street number was parsed, return nil
	if final.StreetNumber == 0 && p.requireStreetNumber {
		return nil
	}
	final.Start, final.Len = identifiers.startLen()

	// If we failed to parse a unit identifier, try to just use the chunk of
	// unparsed string at the beginning of the address (if it exists).
	// See docs on identifiers.unitIdentifier for more info.
	if final.UnitIdentifier == "" {
		unitIdentifier := identifiers.unitIdentifier(addressUpper)
		if unitIdentifier != "" {
			final.Len += final.Start
			final.Start = 0
			final.UnitIdentifier = strings.Trim(unitIdentifier, ", ")
		}
	}

	p.logger.Debugf("final: %#v", *final)
	return final
}

func normaliseUnitType(t string) string {
	t = strings.ToUpper(t)
	switch t {
	case "FLAT", "FLATS", "F":
		t = UnitTypeFlat
	case "LOT", "LOTS":
		t = UnitTypeLot
	case "UNIT", "UNITS", "U":
		t = UnitTypeUnit
	case "VILLA", "VILLAS":
		t = UnitTypeVilla
	}
	return t
}

func (i ParsedAddressIdentifier) String() string {
	return nzAddress.Address{
		UnitType:       i.UnitType,
		UnitIdentifier: i.UnitIdentifier,
		StreetNumber:   int(i.StreetNumber),
		StreetAlpha:    i.StreetAlpha,
	}.Display()
}

// unitIdentifier traverses through the slice of parsed identifiers to find the
// identifier that has the lowest start index (the earliest in the string). If
// it is > 0 (i.e. doesn't start at the beginning), then all characters from the
// beginning of the input string (str) to the first parsed identifier will be
// returned as the unit identifier.
//
// This is basically a catch-all in case we couldn't parse the unit identifier
// some other way. At least one identifier must be in the collection, else the
// method returns an empty string.
//
// Example: "ABC2345 23 SOMETHING STREET"
// If 23 was parsed as the street number and no other identifiers parsed, the
// identifier with the lowest index would be the street number starting at
// position 8 in the string. 0:8 (ABC2345 ) would become the unit identifier.
func (is ParsedAddressIdentifiers) unitIdentifier(str string) string {
	if len(is) == 0 {
		return ""
	}
	min, _ := is.startLen()
	return strings.Trim(str[0:min], " ")
}

// startLen returns the string indexes of the very first and last matched
// characters across all the identifiers in the collection.
func (is ParsedAddressIdentifiers) startLen() (int, int) {
	start := is[0].Start
	len := 0
	for _, i := range is {
		l := i.Start + i.Len
		if l > len {
			len = l
		}
		if i.Start < start {
			start = i.Start
		}
	}
	return start, len
}

// mergeIdentifiers takes the first identifier match in the collection, matches
// the remaining identifiers into it, then returns a single merged result. Will
// remove any conflicting identifiers from the collection.
func (p Parser) mergeIdentifiers(is *ParsedAddressIdentifiers) *ParsedAddressIdentifier {
	identifiers := *is
	final := identifiers[0]
	processedIdentifiers := ParsedAddressIdentifiers{identifiers[0]}
	p.logger.Debug("identifiers: ", *is)
	// Start at 1 to skip the first value, as we've already added it. Don't need
	// to merge an identifier into itself.
	for i := 1; i < len(identifiers); i++ {
		err := final.merge(identifiers[i])
		if err != nil {
			p.logger.Debug("could not merge identifiers: ", err)
			continue
		}
		processedIdentifiers = append(processedIdentifiers, identifiers[i])
	}
	*is = processedIdentifiers
	return final
}

// merge merges an identifier source (s) into the a target (i). Will return an
// error if either identifier overlaps, or if the source disagrees with the
// target about any of the identifier fields.
func (i *ParsedAddressIdentifier) merge(s *ParsedAddressIdentifier) error {
	if s.Len == 0 {
		return errors.New("source len is 0")
	}

	if i.Start < s.Start {
		if i.Start+i.Len > s.Start {
			return errors.New("source overlaps target")
		}
	} else if s.Start < i.Start {
		if s.Start+s.Len > i.Start {
			return errors.New("target overlaps source")
		}
	} else {
		return errors.New("source starts at same place as target")
	}

	// First we check to see if the source disagrees with the target
	if i.UnitType != "" && s.UnitType != "" && i.UnitType != s.UnitType {
		return errors.New("target has different non-empty unit type")
	}
	if i.UnitIdentifier != "" && s.UnitIdentifier != "" && i.UnitIdentifier != s.UnitIdentifier {
		return errors.New("target has different non-empty unit identifier")
	}
	if i.StreetNumber != 0 && s.StreetNumber != 0 && i.StreetNumber != s.StreetNumber {
		return errors.New("target has different non-empty street number")
	}
	if i.StreetNumberHigh != 0 && s.StreetNumberHigh != 0 && i.StreetNumberHigh != s.StreetNumberHigh {
		return errors.New("target has different non-empty street number (high)")
	}
	if i.StreetAlpha != "" && s.StreetAlpha != "" && i.StreetAlpha != s.StreetAlpha {
		return errors.New("target has different non-empty street alpha")
	}

	// Now we check for empty values in target, and overwrite them with values from source
	if i.UnitType == "" {
		i.UnitType = s.UnitType
	}
	if i.UnitIdentifier == "" {
		i.UnitIdentifier = s.UnitIdentifier
	}
	if i.StreetNumber == 0 {
		i.StreetNumber = s.StreetNumber
	}
	if i.StreetNumberHigh == 0 {
		i.StreetNumberHigh = s.StreetNumberHigh
	}
	if i.StreetAlpha == "" {
		i.StreetAlpha = s.StreetAlpha
	}
	return nil
}
