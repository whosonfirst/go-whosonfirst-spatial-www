package geo

func IsValidLatitude(lat float64) bool {
	return lat < 90.00 && lat > -90.00
}

func IsValidLongitude(lon float64) bool {
	return lon < 180.00 && lon > -180.00
}
