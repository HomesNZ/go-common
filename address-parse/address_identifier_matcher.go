package addressParse

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	expRestOfAddress          = "([^0-9\\/]+[0-9]*|$)"
	expUnitIdentifier         = "(([^\\/]+)\\/)"
	expUnitIdentifierWithType = "(((U(?:NIT)?|F(?:LAT)?|L(?:OTS?)?|VILLAS?)\\s*([^\\/,\\.]+))[,\\/\\s]*)"
	expStateHighway           = "((?:STATE\\s*[A-Z]+\\s*|S\\s?H\\s*)([0-9]+[A-Z]?)\\s*)"
	expStreetAlpha            = "(([A-Z]{1,2})(\\s|$))"
	expStreetNumber           = "(([0-9]+))[A-Z]{0,2}\\s*"
	expStreetNumberRange      = "(([0-9]+)-([0-9]+)\\s*)"
)

func (p Parser) invalidMatchers() []IdentifierInvalidMatcher {
	return []IdentifierInvalidMatcher{
		regexInvalidMatcher{
			logger:     p.logger,
			Name:       "state highway without street number",
			Expression: regexp.MustCompile("^\\s*" + expStateHighway + expRestOfAddress),
		},
	}
}

func (p Parser) matchers() []IdentifierMatcher {
	return []IdentifierMatcher{
		regexMatcher{
			logger:         p.logger,
			Name:           "unit identifier containing street alpha",
			Expression:     regexp.MustCompile("^(([A-Z]{1,2})\\s([A-Z0-9]+)\\/)"),
			UnitIdentifier: 3,
			StreetAlpha:    2,
			Len:            1,
		},
		regexMatcher{
			logger:         p.logger,
			Name:           "unit identifier with type",
			Expression:     regexp.MustCompile("^" + expUnitIdentifierWithType),
			UnitType:       3,
			UnitIdentifier: 4,
			Len:            1,
		},
		regexMatcher{
			logger:         p.logger,
			Name:           "unit identifier seperated by forward slash",
			Expression:     regexp.MustCompile("^" + expUnitIdentifier),
			UnitIdentifier: 2,
			Len:            1,
		},
		regexMatcher{
			logger:           p.logger,
			Name:             "unit identifier with type street number range",
			Expression:       regexp.MustCompile("^(" + expUnitIdentifierWithType + ".*?)" + expStreetNumberRange),
			Start:            1,
			StreetNumber:     7,
			StreetNumberHigh: 8,
			Len:              6,
		},
		regexMatcher{
			logger:      p.logger,
			Name:        "unit identifier with type street alpha",
			Expression:  regexp.MustCompile("^(" + expUnitIdentifierWithType + ".*?)[0-9]+" + expStreetAlpha),
			Start:       1,
			StreetAlpha: 7,
			Len:         6,
		},
		regexMatcher{
			logger:       p.logger,
			Name:         "unit identifier with type street number",
			Expression:   regexp.MustCompile("^(" + expUnitIdentifierWithType + ".*?)" + expStreetNumber),
			Start:        1,
			StreetNumber: 7,
			Len:          6,
		},
		regexMatcher{
			logger:      p.logger,
			Name:        "state highway street alpha",
			Expression:  regexp.MustCompile("^(.*?[0-9]+)" + expStreetAlpha + expStateHighway + expRestOfAddress),
			StreetAlpha: 3,
			Start:       1,
			Len:         2,
		},
		regexMatcher{
			logger:           p.logger,
			Name:             "state highway street number range",
			Expression:       regexp.MustCompile("^(.*?)" + expStreetNumberRange + expStateHighway + expRestOfAddress),
			StreetNumber:     3,
			StreetNumberHigh: 4,
			Start:            1,
			Len:              2,
		},
		regexMatcher{
			logger:       p.logger,
			Name:         "state highway street number",
			Expression:   regexp.MustCompile("^(.*?)(([0-9]+)\\s*)" + expStateHighway + expRestOfAddress),
			StreetNumber: 3,
			Start:        1,
			Len:          2,
		},
		regexMatcher{
			logger:      p.logger,
			Name:        "street alpha",
			Expression:  regexp.MustCompile("^(.*?[0-9]+)" + expStreetAlpha + expRestOfAddress),
			StreetAlpha: 3,
			Start:       1,
			Len:         2,
		},
		regexMatcher{
			logger:           p.logger,
			Name:             "street number range",
			Expression:       regexp.MustCompile("^(.*?)" + expStreetNumberRange + expRestOfAddress),
			StreetNumber:     3,
			StreetNumberHigh: 4,
			Start:            1,
			Len:              2,
		},
		regexMatcher{
			logger:       p.logger,
			Name:         "unit identifier street number",
			Expression:   regexp.MustCompile("^(" + expUnitIdentifier + ".*?)" + expStreetNumber + expRestOfAddress),
			StreetNumber: 5,
			Start:        1,
			Len:          4,
		},
		regexMatcher{
			logger:         p.logger,
			Name:           "unit identifier seperated by forward slash",
			Expression:     regexp.MustCompile("^(([A-Z0-9]+)[^A-Z0-9]([0-9]+)([A-Z]{0,2}))"),
			UnitIdentifier: 2,
			StreetNumber:   3,
			StreetAlpha:    4,
			Len:            1,
		},
		regexMatcher{
			logger:       p.logger,
			Name:         "street number",
			Expression:   regexp.MustCompile("^(.*?)" + expStreetNumber + expRestOfAddress),
			StreetNumber: 3,
			Start:        1,
			Len:          2,
		},
	}
}

// IdentifierMatcher is an interface for matchers to parse individual identifier
// fields. If they parse any values, they should return a
// ParsedAddressIdentifier containing the values parsed and the Start+Len to
// indicate where abouts in the string they were parsed from. If no values are
// parsed, the Parse method should return nil.
type IdentifierMatcher interface {
	Parse(string) *ParsedAddressIdentifier
	MatcherName() string
}

// IdentifierInvalidMatcher is an interface for matchers that are used to
// completely invalidate a string as containing an address identifier. Any
// .Invalid method that returns true will return a failed parse result.
type IdentifierInvalidMatcher interface {
	Invalid(string) bool
	MatcherName() string
}

// regexMatcher is an IdentifierMatcher that uses regular expressions to match
// identifier values.
type regexMatcher struct {
	logger logrus.FieldLogger

	Name       string
	Expression *regexp.Regexp
	// The following identifier fields indicate which identifier fields relate to
	// matched regex groups from the Expression.
	//
	// Example: regexMatcher{
	//   Expression:     regexp.MustCompile("^hello( (world))$"),
	//   UnitIdentifier: 2,
	// }.Parse("hello world")
	// Returns: &ParsedAddressIdentifier{
	//   UnitIdentifier: "world",
	// }
	UnitType         int
	UnitIdentifier   int
	StreetNumber     int
	StreetNumberHigh int
	StreetAlpha      int
	// The Start+Len fields are used to determine the precise portion of the
	// string that has been matched. This allows you the matcher to "take
	// responsibility" for portions of the string, even if they are not part of
	// the fields (e.g. dealing with superfluous special characters).
	//
	// Example: regexMatcher{
	//   Expression:     regexp.MustCompile("^(hello)( (world))$"),
	//   Start:          1,
	//   Len:            2,
	//   UnitIdentifier: 3,
	// }.Parse("hello world")
	// Returns: &ParsedAddressIdentifier{
	//   UnitIdentifier: "world", // value is 3rd matched group
	//   Start:          5,       // value is len(1st matched group)
	//   Len:            6,       // value is the len(2nd matched group)
	//                            // (includes the space between "hello world")
	// }
	Start int
	Len   int
}

// Parse implements IdentifierMatcher
func (m regexMatcher) Parse(str string) *ParsedAddressIdentifier {
	parts := m.Expression.FindStringSubmatch(str)
	m.logger.Debugf("%s : %s", str, m.Expression.String())
	if parts == nil {
		return nil
	}
	m.logger.Debugf("parts: %#v", parts)
	parsed := &ParsedAddressIdentifier{Len: len(parts[m.Len])}
	if m.Start != 0 {
		parsed.Start = len(parts[m.Start])
	}
	if m.UnitType != 0 && len(parts) > m.UnitType {
		parsed.UnitType = normaliseUnitType(parts[m.UnitType])
	}
	if m.UnitIdentifier != 0 && len(parts) > m.UnitIdentifier {
		parsed.UnitIdentifier = strings.ToUpper(parts[m.UnitIdentifier])
	}
	if m.StreetNumber != 0 && len(parts) > m.StreetNumber {
		parsed.StreetNumber, _ = strconv.ParseInt(parts[m.StreetNumber], 10, 64)
	}
	if m.StreetNumber != 0 && len(parts) > m.StreetNumberHigh {
		parsed.StreetNumberHigh, _ = strconv.ParseInt(parts[m.StreetNumberHigh], 10, 64)
	}
	if m.StreetAlpha != 0 && len(parts) > m.StreetAlpha {
		parsed.StreetAlpha = strings.ToUpper(parts[m.StreetAlpha])
	}
	return parsed
}

// MatcherName implements IdentifierMatcher
func (m regexMatcher) MatcherName() string {
	return m.Name
}

// RegexInvalidMatcher is a IdentifierInvalidMatcher that simply returns true
// when the Expression matches the input string.
type regexInvalidMatcher struct {
	logger logrus.FieldLogger

	Name       string
	Expression *regexp.Regexp
}

// Invalid implements IdentifierInvalidMatcher
func (m regexInvalidMatcher) Invalid(str string) bool {
	m.logger.Debugf("%s : %s", str, m.Expression.String())
	return m.Expression.MatchString(str)
}

// MatcherName implements IdentifierInvalidMatcher
func (m regexInvalidMatcher) MatcherName() string {
	return m.Name
}
