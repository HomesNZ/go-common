package addressParse

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/HomesNZ/go-common/nz_address"
)

// ParsedStreet represents a street parsed from a string
type ParsedStreet struct {
	Name              string
	UnabbreviatedName string
	Type              string
	Direction         string
	Start             int
	Len               int
}

var (
	// lastNumberExp is what we use to strip out state highway numbers
	//
	// Examples:
	// - "1203 State Highway 23, Auckland": "23"
	// - "State Highway 23, Auckland":      "23"
	// - "58 Kenya Street, Ngaio":          No match
	lastNumberExp = regexp.MustCompile("([0-9]+)((\\s*-\\s*[0-9]+)?[^0-9]*)$")

	stateHighwayExp = regexp.MustCompile("^" + expStateHighway)
)

const (
	_ = iota
	streetPartName
	streetPartType
	streetPartDirection
)

func (s ParsedStreet) String() string {
	return nzAddress.Address{
		StreetName:      s.Name,
		StreetType:      s.Type,
		StreetDirection: s.Direction,
	}.Display()
}

type streetPart struct {
	Type   int
	String string
}

type streetParts struct {
	Parts []streetPart
}

func (s streetParts) get(part int) string {
	parts := []string{}
	for _, v := range s.Parts {
		if v.Type == part {
			parts = append(parts, v.String)
		}
	}
	return strings.Join(parts, " ")
}

func (s streetParts) has(part int) bool {
	for _, v := range s.Parts {
		if v.Type == part {
			return true
		}
	}
	return false
}

func (s *streetParts) set(partFrom int, partTo int) {
	for k, v := range s.Parts {
		if v.Type == partFrom {
			s.Parts[k].Type = partTo
		}
	}
}

// removeTrailingNameParts finds the last non-name-type part and removes all
// parts after that part. Returns number of characters trimmed including spaces.
func (s *streetParts) removeTrailingNameParts() int {
	var rlen int
	if !s.has(streetPartType) {
		return 0
	}
	after := streetPartType
	if s.has(streetPartDirection) {
		after = streetPartDirection
	}
	remove := false
	for k, v := range s.Parts {
		if v.Type == after {
			remove = true
		} else if remove {
			rlen += len(v.String) + 1 // +1 for space
			if len(s.Parts) > k {
				s.Parts = append(s.Parts[:k], s.Parts[k+1:]...)
			} else {
				s.Parts = s.Parts[:k]
			}
		}
	}
	return rlen
}

func (s *streetParts) add(part string, partType int) {
	s.Parts = append(s.Parts, streetPart{
		String: part,
		Type:   partType,
	})
}

// Street parses a street string or returns an error if it can't be parsed
func (p Parser) Street(street string) *ParsedStreet {
	if street == "" {
		return nil
	}

	parts := streetParts{}
	/*
		idxEnd keeps track of the length of the string we are matching from the
		original address string. Because of abbreviations and spaces etc, the
		resulting strings may total a different length so we need to keep track of
		how many characters we are using to build the street name / type.
		Example:
			Chilton Gr (idxEnd == 10) => Chilton Grove (len(street) == 13)
			Cuba  Street (idxEnd == 12) => Cuba Street (len(street) == 11)

		We start at -1 to account for the lack of an initial space - when you want
		to chop a tree into 10 pieces you make 9 cuts.
	*/
	idxEnd := -1
	streetFull := street
	street = strings.TrimLeft(street, " ,")
	comma := strings.Index(street, ",")
	idxEnd += len(streetFull) - len(street)
	if comma != -1 {
		idxEnd++                 // include the comma in our calculation of street name length
		street = street[0:comma] // only match up to the first comma
	}

	// Handle state highways
	stateHighway := stateHighwayExp.FindStringSubmatch(strings.ToUpper(street))
	if stateHighway != nil {
		parts.add(fmt.Sprintf("STATE HIGHWAY %s", stateHighway[2]), streetPartName)
		idxEnd += len(stateHighway[0])
		street = street[len(stateHighway[0]):]
	}

	rawParts := strings.Split(street, " ")
	numParts := len(parts.Parts)
	for _, v := range rawParts {
		if len(v) == 0 {
			continue
		}
		// account for double spaces
		vUpper := strings.Trim(strings.ToUpper(v), " ")

		// full length of part including split + trimmed spaces
		len := len(v) + 1

		if parts.has(streetPartName) {
			// Street type
			if _, ok := p.streetTypeAbbreviations[vUpper]; ok {
				if parts.has(streetPartType) {
					parts.set(streetPartType, streetPartName)
					parts.set(streetPartDirection, streetPartName)
				}
				idxEnd += len
				parts.add(vUpper, streetPartType)
				continue
			}

			// Street direction
			if parts.has(streetPartType) || stateHighway != nil {
				if _, ok := p.streetDirectionAbbreviations[vUpper]; ok {
					if parts.has(streetPartDirection) {
						parts.set(streetPartDirection, streetPartName)
					}
					idxEnd += len
					parts.add(vUpper, streetPartDirection)
					continue
				}
			}
		}

		// Street name
		idxEnd += len
		parts.add(vUpper, streetPartName)
	}

	idxEnd -= parts.removeTrailingNameParts()

	// If we haven't added any more parts, then we move the cursor forward by 1.
	// This is to account for when state highways are matched so we have a matched
	// street, but no additional parts. This breaks the beginning index of -1.
	if numParts == len(parts.Parts) {
		idxEnd++
	}

	parsed := &ParsedStreet{
		Name:      parts.get(streetPartName),
		Type:      parts.get(streetPartType),
		Direction: parts.get(streetPartDirection),
		Start:     0,
		Len:       idxEnd,
	}
	parsed.UnabbreviatedName = parsed.Name

	// Process abbreviations in street names
	p.AbbreviateStreetNameSuburbCity(&parsed.Name)

	if len(parsed.Type) > 0 {
		parsed.Type = p.streetTypeAbbreviations[parsed.Type]

		if parsed.Name == "THE" {
			parsed.Name += " " + parsed.Type
			parsed.UnabbreviatedName += " " + parsed.Type
			parsed.Type = ""
		}
	} else {
		// we assume the street doesnt have a type.
	}
	if len(parsed.Direction) > 0 {
		parsed.Direction = p.streetDirectionAbbreviations[parsed.Direction]
	}

	p.logger.Debugf(
		"Street name: %s, Type: %s, Direction: %s, Start: %d, Len: %d",
		parsed.Name,
		parsed.Type,
		parsed.Direction,
		parsed.Start,
		parsed.Len,
	)

	return parsed
}

// AbbreviateStreetNameSuburbCity is intended for use on street names, suburbs,
// and cities. These have different abbreviations to street types, e.g. ST as a
// street type would be short for STREET, whereas ST in a street name or suburb
// is going to be short for SAINT.
// There are some exceptions which have been hardcoded. Needs improvement in
// this area but this fixes the immediate issues.
func (p Parser) AbbreviateStreetNameSuburbCity(str *string) {
	rawParts := strings.Split(*str, " ")
	// need at least 2 parts for this method to work properly
	if len(rawParts) < 2 {
		// calling methods expect the string to be converted to uppercase
		*str = strings.ToUpper(*str)
		return
	}

	parts := []string{}
	startsWithSaint := false
	for k, v := range rawParts {
		v = strings.ToUpper(v)
		if longForm, ok := p.streetNameSuburbCityAbbreviations[v]; ok {
			// if the street name/suburb/city starts with ST
			if longForm == "SAINT" && k == 0 {
				startsWithSaint = true
			}
			parts = append(parts, longForm)
			continue
		}
		if v == "SAINT" && k == 0 {
			startsWithSaint = true
		} else if startsWithSaint && k == 1 {
			exceptions := []string{
				"ALBANS",
				"ANDREWS",
				"HELIERS",
				"ARNAUD",
			}
			for _, e := range exceptions {
				if strings.HasPrefix(e, v) {
					parts[0] = "ST"
					break
				}
			}
		}
		parts = append(parts, v)
	}
	*str = strings.Join(parts, " ")
}
