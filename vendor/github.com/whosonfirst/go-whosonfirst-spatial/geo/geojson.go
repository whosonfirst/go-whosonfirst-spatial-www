package geo

import (
	"github.com/paulmach/go.geojson"
	"github.com/skelterjohn/geom"
)

func GeoJSONFeatureContainsCoord(f *geojson.Feature, c *geom.Coord) bool {

	return GeoJSONGeometryContainsCoord(f.Geometry, c)
}

func GeoJSONGeometryContainsCoord(geom *geojson.Geometry, c *geom.Coord) bool {

	if geom.IsMultiPolygon() {
		return GeoJSONMultiPolygonContainsCoord(geom.MultiPolygon, c)
	}

	if geom.IsPolygon() {
		return GeoJSONPolygonContainsCoord(geom.Polygon, c)
	}

	return false
}

func GeoJSONMultiPolygonContainsCoord(multi [][][][]float64, c *geom.Coord) bool {

	for _, poly := range multi {

		if GeoJSONPolygonContainsCoord(poly, c) {
			return true
		}
	}

	return false
}

func GeoJSONPolygonContainsCoord(poly [][][]float64, c *geom.Coord) bool {

	count := len(poly)

	if count == 0 {
		return false
	}

	// exterior ring

	exterior_ring := poly[0]

	if !GeoJSONRingContainsCoord(exterior_ring, c) {
		return false
	}

	// interior rings

	if count > 1 {

		for _, interior_ring := range poly {

			if GeoJSONRingContainsCoord(interior_ring, c) {
				return false
			}
		}
	}

	return true
}

func GeoJSONRingContainsCoord(ring [][]float64, c *geom.Coord) bool {

	polygon := geom.Polygon{}

	for _, pt := range ring {
		polygon.AddVertex(geom.Coord{X: pt[0], Y: pt[1]})
	}

	return polygon.ContainsCoord(*c)
}
