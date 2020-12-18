package flags

import (
	"errors"
	"flag"
	"github.com/sfomuseum/go-flags/multi"
)

func AppendQueryFlags(fs *flag.FlagSet) error {

	fs.Float64("latitude", 0.0, "A valid latitude.")
	fs.Float64("longitude", 0.0, "A valid longitude.")

	fs.String("geometries", "all", "Valid options are: all, alt, default.")

	var props multi.MultiString
	fs.Var(&props, "properties", "One or more Who's On First properties to append to each result.")

	var pts multi.MultiString
	fs.Var(&pts, "placetype", "One or more place types to filter results by.")

	var alt_geoms multi.MultiString
	fs.Var(&alt_geoms, "alternate-geometry", "One or more alternate geometry labels (wof:alt_label) values to filter results by.")

	var is_current multi.MultiInt64
	fs.Var(&is_current, "is-current", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_ceased multi.MultiInt64
	fs.Var(&is_ceased, "is-ceased", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_deprecated multi.MultiInt64
	fs.Var(&is_deprecated, "is-deprecated", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_superseded multi.MultiInt64
	fs.Var(&is_superseded, "is-superseded", "One or more existential flags (-1, 0, 1) to filter results by.")

	var is_superseding multi.MultiInt64
	fs.Var(&is_superseding, "is-superseding", "One or more existential flags (-1, 0, 1) to filter results by.")

	return nil
}

func ValidateQueryFlags(fs *flag.FlagSet) error {
	// need to update return Multi wahwah in lookup.go
	return errors.New("Not implemented")
}
