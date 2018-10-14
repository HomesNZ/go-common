package nzAddress

import (
	"strconv"
	"strings"
)

// Address is a complete New Zealand address with all the required fields to
// format into a human-readable address string.
type Address struct {
	// Identifier is for a preformatted unit type/identifier, street number, and
	// street alpha. You can use this instead of the separate address identifier
	// fields.
	Identifier string `json:"identifier"`

	UnitType         string `json:"unit_type"`
	UnitIdentifier   string `json:"unit_identifier"`
	StreetNumber     int    `json:"street_number"`
	StreetNumberHigh int    `json:"street_number_high"`
	StreetAlpha      string `json:"street_alpha"`

	StreetName      string `json:"street_name"`
	StreetType      string `json:"street_type"`
	StreetDirection string `json:"street_direction"`
	Suburb          string `json:"suburb"`
	City            string `json:"city"`

	BuildingName string `json:"building_name"`
	Floor        string `json:"floor"`
	RDNumber     string `json:"rd_number"`
	Postcode     int    `json:"postcode"`
}

// Street returns the formatted street name + type and direction
func (a Address) Street() string {
	parts := []string{a.StreetName}
	if a.StreetType != "" {
		parts = append(parts, a.StreetType)
	}
	if a.StreetDirection != "" {
		parts = append(parts, a.StreetDirection)
	}
	return strings.Trim(strings.Join(parts, " "), " ")
}

func titleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// Display formats an address into a valid display address
func (a Address) Display() string {
	address := []string{}

	if a.BuildingName != "" {
		unitBuildingName := titleCase(a.BuildingName)
		if a.UnitIdentifier != "" {
			unitBuildingName = strings.ToUpper(a.UnitIdentifier) + " " + unitBuildingName
			if a.UnitType != "" {
				unitBuildingName = titleCase(a.UnitType) + " " + unitBuildingName
			}
		}
		address = append(address, unitBuildingName)
	}

	var identifierStreet string
	street := titleCase(a.Street())
	if a.Identifier != "" {
		identifierStreet += a.Identifier

		if street != "" {
			identifierStreet += " "
		}
	} else if a.StreetNumber != 0 {
		if a.UnitIdentifier != "" && a.BuildingName == "" {
			if a.UnitType != "" {
				identifierStreet += titleCase(a.UnitType) + " "
			}
			identifierStreet += strings.ToUpper(a.UnitIdentifier) + "/"
		}

		identifierStreet += strconv.Itoa(a.StreetNumber) + strings.ToUpper(a.StreetAlpha)
		if a.StreetNumberHigh != 0 {
			identifierStreet += "-" + strconv.Itoa(a.StreetNumberHigh)
		}

		if street != "" {
			identifierStreet += " "
		}
	}
	identifierStreet += street
	if identifierStreet != "" {
		address = append(address, identifierStreet)
	}

	suburb := titleCase(a.Suburb)
	city := titleCase(a.City)

	if suburb != "" && city != suburb {
		address = append(address, suburb)
	}
	if a.RDNumber != "" {
		address = append(address, "RD "+strings.ToUpper(a.RDNumber))
	}
	if city != "" {
		address = append(address, city)
	}
	return strings.Trim(strings.Join(address, ", "), " ")
}

// DisplayWithPostcode returns a display address with postcode appended to the
// end.
func (a Address) DisplayWithPostcode() string {
	return a.Display() + " " + strconv.Itoa(a.Postcode)
}
