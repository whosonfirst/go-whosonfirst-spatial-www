package geo

import (
	"github.com/skelterjohn/geom"
)

func MultiPolygonContainsCoord(multi [][][][]float64, c *geom.Coord) bool {

	for _, poly := range multi {

		if PolygonContainsCoord(poly, c) {
			return true
		}
	}

	return false
}

func PolygonContainsCoord(poly [][][]float64, c *geom.Coord) bool {

	count := len(poly)

	if count == 0 {
		return false
	}

	// exterior ring

	exterior_ring := poly[0]

	if !RingContainsCoord(exterior_ring, c) {
		return false
	}

	// interior rings

	if count > 1 {

		for _, interior_ring := range poly[1:] {

			if RingContainsCoord(interior_ring, c) {
				return false
			}
		}
	}

	return true
}

func RingContainsCoord(ring [][]float64, c *geom.Coord) bool {

	polygon := geom.Polygon{}

	for _, pt := range ring {
		polygon.AddVertex(geom.Coord{X: pt[0], Y: pt[1]})
	}

	return polygon.ContainsCoord(*c)
}
