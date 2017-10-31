package util

import (
	"fmt"
	"regexp"
	"strconv"
	"math"
)

var (
	buildingAgeRegexp = regexp.MustCompile("..([\\d\\d\\d]+).*") // This matches the building age part of the CategoryCode
)

// ParseBool attempts to take a string value and return a bool. It wraps strconv.ParseBool adding "yes", "YES", "Yes", "Y", "y" as TRUE values,
// and "no", "NO", "No", "N", "n", "NULL", and empty string as FALSE values.
func ParseBool(str string) (value bool, err error) {
	switch str {
	case "yes", "YES", "Yes", "Y", "y":
		return true, nil
	case "no", "NO", "No", "N", "n", "NULL", "":
		return false, nil
	default:
		ival, err := strconv.ParseBool(str)
		if err != nil {
			return false, fmt.Errorf("Unable to parse <%s> as bool", str)
		}
		return ival, nil
	}
}

// ForceParseBool is the same as ParseBool except returns false if it fails to
// parse instead of returning an error.
func ForceParseBool(str string) bool {
	parsed, _ := ParseBool(str)
	return parsed
}

// IsEmpty returns true if a string is considered to be "empty", otherwise false.
func IsEmpty(str string) bool {
	switch str {
	case "null", "NULL", "", "0":
		return true
	}
	return false
}

func ToSquareMeters(val float64, unit string) (float64, error) {
	// Return early if we have zero value as theres nothing to convert.
	if val == 0 {
		return 0, nil
	}

	var s float64

	switch unit {
	// square is an Imperial unit of area that is used in the United States construction industry, and was historically
	// used in Australia.
	case "square":
		s = val * 9.290304
	case "squareMeter":
		s = val
	case "a", "acre", "acres":
		s = val * 4046.86
	case "h", "hectare", "hectares":
		s = val * 10000.0
	default:
		return 0, fmt.Errorf("Cannot convert to square meters, unsupported unit `%s`", unit)
	}

	return math.Ceil(s - 0.5), nil
}

// GetBuildingAgeFromCategoryCode retrieves a CategoryCode string, eg: RD201A, and returns the '201' bit, which is the
// decade built indicator. If there is no match (ie, no 3 numbers), it returns ""
//
// This method is required here, but could be in a more general place.
func GetBuildingAgeFromCategoryCode(categoryCode string) string {
	buildingAgeIndicator := ""
	matches := buildingAgeRegexp.FindStringSubmatch(categoryCode)
	if len(matches) > 1 {
		buildingAgeIndicator = matches[1]
	}
	return buildingAgeIndicator
}
