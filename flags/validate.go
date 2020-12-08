package flags

import (
	"errors"
	"flag"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
)

func ValidateWWWFlags(fs *flag.FlagSet) error {

	enable_www, err := spatial_flags.BoolVar(fs, "enable-www")

	if err != nil {
		return err
	}

	if !enable_www {
		return nil
	}

	log.Println("-enable-www flag is true causing the following flags to also be true: -enable-geojson -enable-candidates -enable-properties")

	fs.Set("enable-geojson", "true")
	fs.Set("enable-properties", "true")
	fs.Set("enable-candidates", "true")

	properties_reader_uri, err := StringVar(fs, "properties-reader-uri")

	if err != nil {
		return err
	}

	if properties_reader_uri == "" {
		return errors.New("Invalid or missing -properties-reader-uri flag")
	}

	init_lat, err := spatial_flags.Float64Var(fs, "initial-latitude")

	if err != nil {
		return err
	}

	if !geo.IsValidLatitude(init_lat) {
		return errors.New("Invalid latitude")
	}

	init_lon, err := spatial_flags.Float64Var(fs, "initial-longitude")

	if err != nil {
		return err
	}

	if !geo.IsValidLongitude(init_lon) {
		return errors.New("Invalid longitude")
	}

	init_zoom, err := spatial_flags.IntVar(fs, "initial-zoom")

	if err != nil {
		return err
	}

	if init_zoom < 1 {
		return errors.New("Invalid zoom")
	}

	return nil
}
