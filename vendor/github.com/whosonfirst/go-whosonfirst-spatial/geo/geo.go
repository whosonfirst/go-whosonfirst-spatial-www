package geo

import (
	"errors"
	"github.com/skelterjohn/geom"
)

func IsValidMinLatitude(lat float64) bool {
	return lat >= -90.00
}

func IsValidMaxLatitude(lat float64) bool {
	return lat <= 90.00
}

func IsValidLatitude(lat float64) bool {
	return IsValidMinLatitude(lat) && IsValidMaxLatitude(lat)
}

func IsValidMinLongitude(lon float64) bool {
	return lon >= -180.00
}

func IsValidMaxLongitude(lon float64) bool {
	return lon <= 180.00
}

func IsValidLongitude(lon float64) bool {
	return IsValidMinLongitude(lon) && IsValidMaxLongitude(lon)
}

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

func NewBoundingBox(minx float64, miny float64, maxx float64, maxy float64) (*geom.Rect, error) {

	if !IsValidLongitude(minx) {
		return nil, errors.New("Invalid min longitude")
	}

	if !IsValidLatitude(miny) {
		return nil, errors.New("Invalid min latitude")
	}

	if !IsValidLongitude(maxx) {
		return nil, errors.New("Invalid max longitude")
	}

	if !IsValidLatitude(maxy) {
		return nil, errors.New("Invalid max latitude")
	}

	if minx > maxx {
		return nil, errors.New("Min lon is greater than max lon")
	}

	if minx > maxx {
		return nil, errors.New("Min latitude is greater than max latitude")
	}

	min_coord, err := NewCoordinate(minx, miny)

	if err != nil {
		return nil, err
	}

	max_coord, err := NewCoordinate(maxx, maxy)

	if err != nil {
		return nil, err
	}

	rect := &geom.Rect{
		Min: *min_coord,
		Max: *max_coord,
	}

	return rect, nil
}
