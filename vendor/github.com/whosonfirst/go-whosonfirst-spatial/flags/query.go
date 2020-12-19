package flags

import (
	"errors"
	"flag"
	"github.com/sfomuseum/go-flags/multi"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
)

func AppendQueryFlags(fs *flag.FlagSet) error {

	fs.Float64("latitude", 0.0, "A valid latitude.")
	fs.Float64("longitude", 0.0, "A valid longitude.")

	fs.String("geometries", "all", "Valid options are: all, alt, default.")

	var props multi.MultiString
	fs.Var(&props, "properties", "One or more Who's On First properties to append to each result.")

	var placetypes multi.MultiString
	fs.Var(&placetypes, "placetype", "One or more place types to filter results by.")

	var alt_geoms multi.MultiString
	fs.Var(&alt_geoms, "alternate-geometry", "One or more alternate geometry labels (wof:alt_label) values to filter results by.")

	var is_current multi.MultiString
	fs.Var(&is_current, "is-current", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_ceased multi.MultiString
	fs.Var(&is_ceased, "is-ceased", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_deprecated multi.MultiString
	fs.Var(&is_deprecated, "is-deprecated", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_superseded multi.MultiString
	fs.Var(&is_superseded, "is-superseded", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_superseding multi.MultiString
	fs.Var(&is_superseding, "is-superseding", "One or more existential flags (-1, 0, 1) to filter results by.")

	return nil
}

func ValidateQueryFlags(fs *flag.FlagSet) error {

	lat, err := Float64Var(fs, "latitude")

	if err != nil {
		return err
	}

	lon, err := Float64Var(fs, "longitude")

	if err != nil {
		return err
	}

	if !geo.IsValidLatitude(lat) {
		return errors.New("Invalid latitude")
	}

	if !geo.IsValidLongitude(lon) {
		return errors.New("Invalid longitude")
	}

	_, err = StringVar(fs, "geometries")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "alternate-geometry")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "is-current")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "is-ceased")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "is-deprecated")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "is-superseding")

	if err != nil {
		return err
	}

	_, err = MultiStringVar(fs, "is-superseded")

	if err != nil {
		return err
	}

	return nil
}
