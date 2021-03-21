package flags

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
)

func CommonFlags() (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("common")

	// spatial databases

	available_databases := database.Schemes()
	desc_databases := fmt.Sprintf("Valid options are: %s", available_databases)

	fs.String(SPATIAL_DATABASE_URI, "rtree://", desc_databases)

	fs.Bool("is-wof", true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	// property readers

	fs.Bool(ENABLE_PROPERTIES, false, "Enable support for 'properties' parameters in queries.")
	fs.Bool(INDEX_PROPERTIES, false, "Index properties reader.")

	available_property_readers := properties.Schemes()
	desc_property_readers := fmt.Sprintf("Valid options are: %s", available_property_readers)

	fs.String(PROPERTIES_READER_URI, "", desc_property_readers)

	fs.Bool(ENABLE_CUSTOM_PLACETYPES, false, "...")
	fs.String(CUSTOM_PLACETYPES_SOURCE, "", "...")
	fs.String(CUSTOM_PLACETYPES, "", "...")

	// this is invoked/used in app/indexer.go but for the life of me I can't
	// figure out how to make the code in flags/exclude.go implement the
	// correct inferface wah wah so that flag.Lookup("exclude").Value returns
	// something we can loop over... so instead we just strings.Split() on
	// flag.Lookup("exclude").String() which is dumb but works...
	// (20180301/thisisaaronland)

	var exclude Exclude
	fs.Var(&exclude, EXCLUDE, "Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.")

	fs.Bool(VERBOSE, false, "Be chatty.")

	return fs, nil
}

func ValidateCommonFlags(fs *flag.FlagSet) error {

	spatial_database_uri, err := lookup.StringVar(fs, SPATIAL_DATABASE_URI)

	if err != nil {
		return err
	}

	if spatial_database_uri == "" {
		return fmt.Errorf("Invalid or missing -%s flag", SPATIAL_DATABASE_URI)
	}

	enable_properties, err := lookup.BoolVar(fs, ENABLE_PROPERTIES)

	if err != nil {
		return err
	}

	if enable_properties {

		properties_reader_uri, err := lookup.StringVar(fs, PROPERTIES_READER_URI)

		if err != nil {
			return err
		}

		if properties_reader_uri == "" {
			return fmt.Errorf("Invalid or missing -%s flag", PROPERTIES_READER_URI)
		}

		_, err = lookup.BoolVar(fs, INDEX_PROPERTIES)

		if err != nil {
			return err
		}

	}

	return nil
}
