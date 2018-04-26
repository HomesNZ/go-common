package geo

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/HomesNZ/go-common/extmath"
)

// Point represents a lat/long
type Point struct {
	Lat  float64
	Long float64
	SRID int

	// A flag that stipulates this point must include the longitude twice when
	// marshalling as JSON, once as "long" and once as "lon" (deprecated).
	marshalLongAsLon bool
}

// NewPoint creates and returns a new Point
func NewPoint(x float64, y float64, srid int) *Point {
	return &Point{
		Long: x,
		Lat:  y,
		SRID: srid,
	}
}

// NewPointLongAsLon returns a new point with marshalLongAsLon set to true
func NewPointLongAsLon(x float64, y float64, srid int) *Point {
	return &Point{
		Long:             x,
		Lat:              y,
		SRID:             srid,
		marshalLongAsLon: true,
	}
}

// IsNull returns a boolean indicating whether the point is considered null.
// Currently returns true when the SRID has not been set.
func (p Point) IsNull() bool {
	return p.SRID == 0 || p.Lat == 0 || p.Long == 0
}

// MarshalJSON implements the json.Marshaler interface
func (p Point) MarshalJSON() ([]byte, error) {
	if p.IsNull() {
		return []byte("null"), nil
	}
	if p.marshalLongAsLon {
		return json.Marshal(map[string]float64{
			"lat":  p.Lat,
			"long": p.Long,
			"lon":  p.Long,
		})
	}
	return json.Marshal(map[string]float64{
		"lat":  p.Lat,
		"long": p.Long,
	})
}

func (p *Point) UnmarshalJSON(b []byte) error {
	var v struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
		SRID int     `json:"srid"`
	}
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	p.Lat = v.Lat
	p.Long = v.Long
	switch v.SRID {
	case WGS84, NZTM:
		p.SRID = v.SRID
	default:
		p.SRID = WGS84
	}
	return nil
}

func (p Point) String() string {
	if p.IsNull() {
		return "NULL"
	}
	point := fmt.Sprintf(
		"ST_GeometryFromText('POINT(%s %s)', %d)",
		strconv.FormatFloat(p.Long, 'f', -1, 64),
		strconv.FormatFloat(p.Lat, 'f', -1, 64),
		p.SRID,
	)
	if p.SRID != WGS84 {
		point = fmt.Sprintf("ST_Transform(%s,%d)", point, WGS84)
	}
	return point
}

// Compare returns true if the SRIDs of the given points are the same, and they
// are located within "diff" distance of each other.
// This is not a safe way to compare a lat/long in general, but for the NZ
// coordinates we are dealing with it should work just fine.
func (p Point) Compare(c Point, diff float64) bool {
	if p.SRID != c.SRID {
		return false
	}
	return math.Abs(p.Lat-c.Lat) < diff && math.Abs(p.Long-c.Long) < diff
}

// Distance function returns the distance (in meters) between two points of
//     a given longitude and latitude relatively accurately (using a spherical
//     approximation of the Earth) through the Haversin Distance Formula for
//     great arc distance on a sphere with accuracy for small distances
//
// point coordinates are supplied in degrees and converted into rad. in the func
//
// distance returned is METERS!!!!!!
// http://en.wikipedia.org/wiki/Haversine_formula
//
// Shamelessly stolen from this gist:
// https://gist.github.com/cdipaolo/d3f8db3848278b49db68
// - dylannz
func (p *Point) Distance(p2 *Point) float64 {
	// convert to radians
	la1 := p.Lat * math.Pi / 180
	lo1 := p.Long * math.Pi / 180
	la2 := p2.Lat * math.Pi / 180
	lo2 := p2.Long * math.Pi / 180

	// must cast radius as float to multiply later
	r := 6378100.0 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// Round rounds the lat/long of a point to requested significant figures sf.
func (p *Point) Round(sf int) {
	p.Lat = extmath.Round(p.Lat, sf)
	p.Long = extmath.Round(p.Long, sf)
}
