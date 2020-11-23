package geom

import (
	"fmt"

	geojson "github.com/paulmach/go.geojson"
	"github.com/pkg/errors"
	"github.com/twpayne/go-polyline"
)

func GeometryToCoordinates(geom *geojson.Geometry) [][]float64 {
	switch {
	case geom.IsCollection():
		return CollectionToCoordinates(geom)
	case geom.IsMultiPolygon():
		return MultiPolygonToCoordinates(geom.MultiPolygon)
	case geom.IsPolygon():
		return PolygonToCoordinates(geom.Polygon)
	case geom.IsMultiLineString():
		return PolygonToCoordinates(geom.MultiLineString)
	case geom.IsLineString():
		return geom.LineString
	case geom.IsMultiPoint():
		return geom.MultiPoint
	case geom.IsPoint():
		return [][]float64{[]float64{geom.Point[0], geom.Point[1]}}
	default:
		return [][]float64{}
	}
}

func CollectionToCoordinates(collection *geojson.Geometry) [][]float64 {
	var coordinates [][]float64
	geoms := collection.Geometries
	for _, geom := range geoms {
		coordinates = append(coordinates, GeometryToCoordinates(geom)...)
	}
	return coordinates
}

func MultiPolygonToCoordinates(multigon [][][][]float64) [][]float64 {
	var coordinates [][]float64
	for _, polygon := range multigon {
		coordinates = append(coordinates, PolygonToCoordinates(polygon)...)
	}
	return coordinates
}

func PolylinesToGeoms(polys [][]byte) ([]*geojson.Geometry, error) {
	var geoms []*geojson.Geometry
	for i, poly := range polys {
		geom, err := PolylineToGeom(poly)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("encountered on polyine at idx %d", i))
		}
		geoms = append(geoms, geom)
	}
	return geoms, nil
}

func PolylineToGeom(poly []byte) (*geojson.Geometry, error) {
	coords, _, err := polyline.DecodeCoords(poly)
	if err != nil {
		return nil, errors.Wrap(err, "decoing polyline")
	}

	return geojson.NewPolygonGeometry([][][]float64{
		coords,
	}), err
}

func GeomToEncodedPolyline(geom *geojson.Geometry) [][]byte {
	polylines := [][]byte{}
	switch {
	case geom.IsPolygon():
		for _, line := range geom.Polygon {
			polylines = append(polylines, polyline.EncodeCoords(line))
		}
	case geom.IsMultiPolygon():
		for _, polygon := range geom.MultiPolygon {
			for _, line := range polygon {
				polylines = append(polylines, polyline.EncodeCoords(line))
			}
		}
	}

	return polylines
}

func PolygonToCoordinates(polygon [][][]float64) [][]float64 {
	var coordinates [][]float64
	for _, coord := range polygon {
		coordinates = append(coordinates, coord...)
	}
	return coordinates
}
