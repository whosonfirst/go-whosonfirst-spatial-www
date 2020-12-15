package geo

import (
	"errors"
	"github.com/skelterjohn/geom"
)

func NewCoordinate(x float64, y float64) (*geom.Coord, error) {

	if !IsValidLatitude(y) {
		return nil, errors.New("Invalid latitude")
	}

	if !IsValidLongitude(y) {
		return nil, errors.New("Invalid longitude")
	}

	coord := &geom.Coord{
		X: x,
		Y: y,
	}

	return coord, nil
}
