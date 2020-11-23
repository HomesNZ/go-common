package addressParse

import (
	"github.com/sirupsen/logrus"
)

type Option func(*Parser) error

type Parser struct {
	logger logrus.FieldLogger

	requireStreetNumber bool

	// streetTypeAbbreviations is a map of street type abbreviation => full
	streetTypeAbbreviations map[string]string

	// streetDirectionAbbreviations is a map of street type abbreviation => full
	streetDirectionAbbreviations map[string]string

	// streetNameSuburbCityAbbreviations is a map of street name, suburb, and city
	// abbreviations => full. We don't have a GeoPAF data set for this, hence why
	// it is hardcoded.
	streetNameSuburbCityAbbreviations map[string]string
}

/*
The parseAddress package differs slightly from the old geocoder package. It's
set up so you can load data in from anywhere you like, rather than only being
able to load data from the geocoder.

New usage example below:

```
import "github.com/HomesNZ/geocoder/client"

c := client.New()
abbreviations, err := c.Abbreviations()
if err != nil {
	...
}

parser, err := addressParse.New(
  addressParse.StreetTypeAbbreviations(abbreviations.StreetTypes),
  addressParse.StreetDirectionAbbreviations(abbreviations.StreetDirections),
)
if err != nil {
	...
}
```
*/

func New(opts ...Option) (*Parser, error) {
	p := &Parser{
		requireStreetNumber: true,
		streetNameSuburbCityAbbreviations: map[string]string{
			"PT": "POINT",
			"MT": "MOUNT",
			"ST": "SAINT",
		},
	}
	for _, opt := range opts {
		err := opt(p)
		if err != nil {
			return nil, err
		}
	}
	if p.streetTypeAbbreviations == nil {
		p.streetTypeAbbreviations = map[string]string{}
	}
	if p.streetDirectionAbbreviations == nil {
		p.streetDirectionAbbreviations = map[string]string{}
	}
	if p.logger == nil {
		p.logger = logrus.StandardLogger()
	}
	return p, nil
}

// RequireStreetNumber forces the address identifier parser to return nil if no
// street number was parsed. If set to false, parser will return whatever it has
// managed to parse, regardless of if it contains a street number.
func RequireStreetNumber(v bool) Option {
	return func(p *Parser) error {
		p.requireStreetNumber = v
		return nil
	}
}

func StreetTypeAbbreviations(v map[string]string) Option {
	return func(p *Parser) error {
		p.streetTypeAbbreviations = v
		return nil
	}
}

func StreetDirectionAbbreviations(v map[string]string) Option {
	return func(p *Parser) error {
		p.streetDirectionAbbreviations = v
		return nil
	}
}

func StreetNameSuburbCityAbbreviations(v map[string]string) Option {
	return func(p *Parser) error {
		p.streetNameSuburbCityAbbreviations = v
		return nil
	}
}
