package geom

import (
	"strconv"
	"testing"

	geojson "github.com/paulmach/go.geojson"
)

func TestPolygonToEncodedPolyline(t *testing.T) {
	t.Parallel()

	tests := []struct {
		polygon          *geojson.Geometry
		expectedPolyline [][]byte
	}{
		{
			geojson.NewPolygonGeometry([][][]float64{
				{
					{-41.22, 174.86},
					{-41.23, 174.86},
					{-41.23, 174.85},
					{-41.22, 174.85},
					{-41.22, 174.86},
				},
			}),
			[][]byte{[]byte("~wqzF_jgj`@n}@??n}@o}@??o}@")},
		},
		{
			geojson.NewPolygonGeometry([][][]float64{
				{
					{-41.22, 174.86},
					{-41.23, 174.86},
					{-41.23, 174.85},
					{-41.22, 174.85},
					{-41.22, 174.86},
				},
				{
					{-41.12, 174.96},
					{-41.13, 174.96},
					{-41.13, 174.95},
					{-41.12, 174.95},
					{-41.12, 174.96},
				},
			}),
			[][]byte{[]byte("~wqzF_jgj`@n}@??n}@o}@??o}@"), []byte("~f~yF_{zj`@n}@??n}@o}@??o}@")},
		},
	}

	for _, test := range tests {
		encoded := GeomToEncodedPolyline(test.polygon)

		if len(test.expectedPolyline) != len(encoded) {
			t.Errorf("expected array length [%d] did not match actual array length [%d]", len(test.expectedPolyline), len(encoded))
			return
		}
		for i := range encoded {
			if len(test.expectedPolyline[i]) != len(encoded[i]) {
				t.Errorf("expected %s but got %s", test.expectedPolyline[i], encoded[i])
				return
			}
			for j := range encoded[i] {
				if encoded[i][j] != test.expectedPolyline[i][j] {
					t.Errorf("expected [%s] but got [%s]", test.expectedPolyline[i], encoded[i])
					return
				}
			}

		}
	}
}

func TestPolyLinesToGeoms(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expectedPolygons []*geojson.Geometry
		polyline         [][]byte
	}{
		{
			[]*geojson.Geometry{
				geojson.NewPolygonGeometry([][][]float64{
					{
						{-41.22, 174.86},
						{-41.23, 174.86},
						{-41.23, 174.85},
						{-41.22, 174.85},
						{-41.22, 174.86},
					},
				}),
			},
			[][]byte{[]byte("~wqzF_jgj`@n}@??n}@o}@??o}@")},
		},
		{
			[]*geojson.Geometry{
				geojson.NewPolygonGeometry([][][]float64{
					{
						{-41.22, 174.86},
						{-41.23, 174.86},
						{-41.23, 174.85},
						{-41.22, 174.85},
						{-41.22, 174.86},
					},
				}),
				geojson.NewPolygonGeometry([][][]float64{
					{
						{-41.12, 174.96},
						{-41.13, 174.96},
						{-41.13, 174.95},
						{-41.12, 174.95},
						{-41.12, 174.96},
					},
				}),
			},
			[][]byte{[]byte("~wqzF_jgj`@n}@??n}@o}@??o}@"), []byte("~f~yF_{zj`@n}@??n}@o}@??o}@")},
		},
	}

	for _, test := range tests {
		geoms, err := PolylinesToGeoms(test.polyline)
		if err != nil {
			t.Errorf("encountered unexpected error %s", err)
			return
		}
		if len(geoms) != len(test.expectedPolygons) {
			t.Errorf("expected array length [%d] did not match actual array length [%d]", len(test.expectedPolygons), len(geoms))
			return
		}
		for i := range geoms {
			if len(geoms[i].Polygon) != len(test.expectedPolygons[i].Polygon) {
				t.Errorf("expected polygon length [%d] did not match actual polygon length [%d]", len(test.expectedPolygons[i].Polygon), len(geoms[i].Polygon))
				return
			}
			for j := range geoms[i].Polygon {
				if len(geoms[i].Polygon[j]) != len(test.expectedPolygons[i].Polygon[j]) {
					t.Errorf("expected polygon length [%d] did not match actual polygon length [%d]", len(test.expectedPolygons[i].Polygon[j]), len(geoms[i].Polygon[j]))
					return
				}
				for k := range geoms[i].Polygon[j] {
					actualCoord := geoms[i].Polygon[j][k]
					expectedCoord := test.expectedPolygons[i].Polygon[j][k]
					if strconv.FormatFloat(actualCoord[0], 'f', 5, 64) != strconv.FormatFloat(expectedCoord[0], 'f', 5, 64) ||
						strconv.FormatFloat(actualCoord[1], 'f', 5, 64) != strconv.FormatFloat(expectedCoord[1], 'f', 5, 64) {
						t.Errorf("expected [%#v] but got [%#v]", test.expectedPolygons[i].Polygon, geoms[i].Polygon)
					}
				}
			}
		}
	}
}