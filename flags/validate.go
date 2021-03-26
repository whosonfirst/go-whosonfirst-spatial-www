package flags

import (
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"log"
	"strconv"
	"strings"
)

func ValidateWWWFlags(fs *flag.FlagSet) error {

	enable_www, err := lookup.BoolVar(fs, ENABLE_WWW)

	if err != nil {
		return err
	}

	if !enable_www {
		return nil
	}

	bool_flags := []string{
		ENABLE_CORS,
		ENABLE_GZIP,
		ENABLE_GEOJSON,
	}

	for _, k := range bool_flags {

		_, err := lookup.BoolVar(fs, k)

		if err != nil {
			return err
		}
	}

	init_lat, err := lookup.Float64Var(fs, INITIAL_LATITUDE)

	if err != nil {
		return err
	}

	if !geo.IsValidLatitude(init_lat) {
		return errors.New("Invalid latitude")
	}

	init_lon, err := lookup.Float64Var(fs, INITIAL_LONGITUDE)

	if err != nil {
		return err
	}

	if !geo.IsValidLongitude(init_lon) {
		return errors.New("Invalid longitude")
	}

	init_zoom, err := lookup.IntVar(fs, INITIAL_ZOOM)

	if err != nil {
		return err
	}

	if init_zoom < 1 {
		return errors.New("Invalid zoom")
	}

	max_bounds, err := lookup.StringVar(fs, MAX_BOUNDS)

	if max_bounds != "" {

		bounds := strings.Split(max_bounds, ",")

		if len(bounds) != 4 {
			return errors.New("Invalid max bounds")
		}

		minx, err := strconv.ParseFloat(bounds[0], 64)

		if err != nil {
			return fmt.Errorf("Invalid minx, %v", err)
		}

		if !geo.IsValidLongitude(minx) {
			return errors.New("Invalid longitude, minx")
		}

		miny, err := strconv.ParseFloat(bounds[1], 64)

		if err != nil {
			return fmt.Errorf("Invalid miny, %v", err)
		}

		if !geo.IsValidLatitude(miny) {
			return errors.New("Invalid latitude, miny")
		}

		maxx, err := strconv.ParseFloat(bounds[2], 64)

		if err != nil {
			return fmt.Errorf("Invalid maxx, %v", err)
		}

		if !geo.IsValidLongitude(maxx) {
			return errors.New("Invalid longitude, maxx")
		}

		maxy, err := strconv.ParseFloat(bounds[3], 64)

		if err != nil {
			return fmt.Errorf("Invalid maxy, %v", err)
		}

		if !geo.IsValidLatitude(maxy) {
			return errors.New("Invalid latitude, maxy")
		}
	}

	path_flags := []string{
		PATH_PREFIX,
		PATH_API,
		PATH_DATA,
		PATH_PING,
		PATH_PIP,
	}

	for _, fl := range path_flags {

		_, err := lookup.StringVar(fs, fl)

		if err != nil {
			return err
		}
	}

	enable_tangram, err := lookup.BoolVar(fs, ENABLE_TANGRAM)

	if err != nil {
		return err
	}

	if enable_tangram {

		nz_keys := []string{
			NEXTZEN_APIKEY,
			NEXTZEN_STYLE_URL,
			NEXTZEN_TILE_URL,
		}

		for _, k := range nz_keys {

			v, err := lookup.StringVar(fs, k)

			if err != nil {
				return err
			}

			if v == "" {
				log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", k)
			}
		}

	} else {

		v, err := lookup.StringVar(fs, LEAFLET_TILE_URL)

		if err != nil {
			return err
		}

		if v == "" {
			log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", LEAFLET_TILE_URL)
		}
	}

	return nil
}
