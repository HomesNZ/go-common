package addressParse

import (
	"regexp"
	"strings"
)

var (
	postcodeRegexp = regexp.MustCompile("^([A-Za-z ]+)([0-9]*)")
)

// SuburbAndCity parses a suburb and/or city from the address string. This
// method only works when supplied with the remaining characters of an address
// string after address identifier and street have been removed.
func (p Parser) SuburbAndCity(suburbAndCity string) (string, string, string) {
	suburb := ""
	city := ""
	postcode := ""

	// This is a very dodgy way to do this, but should work for most cases.
	rawParts := strings.Split(suburbAndCity, ",")
	i := 0
	for _, v := range rawParts {
		i++
		switch i {
		case 1:
			suburb = strings.Trim(v, " ")
		case 2:
			parts := postcodeRegexp.FindStringSubmatch(v)
			if parts == nil {
				break
			}
			city = strings.Trim(parts[1], " ")
			if len(parts) > 2 {
				postcode = strings.Trim(parts[2], " ")
			}
		default:
			break
		}
	}

	// Process abbreviations for suburbs and cities
	p.AbbreviateStreetNameSuburbCity(&suburb)
	p.AbbreviateStreetNameSuburbCity(&city)

	return suburb, city, postcode
}
