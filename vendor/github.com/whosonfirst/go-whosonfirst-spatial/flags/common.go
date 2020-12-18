package flags

import (
	"errors"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
)

func CommonFlags() (*flag.FlagSet, error) {

	fs := NewFlagSet("common")

	// spatial databases

	available_databases := database.Schemes()
	desc_databases := fmt.Sprintf("Valid options are: %s", available_databases)

	fs.String("spatial-database-uri", "rtree://", desc_databases)

	// property readers

	fs.Bool("enable-properties", false, "Enable support for 'properties' parameters in queries.")

	available_property_readers := properties.Schemes()
	desc_property_readers := fmt.Sprintf("Valid options are: %s", available_property_readers)

	fs.String("properties-reader-uri", "", desc_property_readers)

	fs.Bool("is-wof", true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	fs.Bool("enable-custom-placetypes", false, "...")
	fs.String("custom-placetypes-source", "", "...")
	fs.String("custom-placetypes", "", "...")

	// this is invoked/used in app/indexer.go but for the life of me I can't
	// figure out how to make the code in flags/exclude.go implement the
	// correct inferface wah wah so that flag.Lookup("exclude").Value returns
	// something we can loop over... so instead we just strings.Split() on
	// flag.Lookup("exclude").String() which is dumb but works...
	// (20180301/thisisaaronland)

	var exclude Exclude
	fs.Var(&exclude, "exclude", "Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.")

	fs.Bool("setenv", false, "Set flags from environment variables.")
	fs.Bool("verbose", false, "Be chatty.")

	return fs, nil
}

func ValidateCommonFlags(fs *flag.FlagSet) error {

	spatial_database_uri, err := StringVar(fs, "spatial-database-uri")

	if err != nil {
		return err
	}

	if spatial_database_uri == "" {
		return errors.New("Invalid or missing -spatial-database-uri flag")
	}

	enable_properties, err := BoolVar(fs, "enable-properties")

	if err != nil {
		return err
	}

	if enable_properties {

		properties_reader_uri, err := StringVar(fs, "properties-reader-uri")

		if err != nil {
			return err
		}

		if properties_reader_uri == "" {
			return errors.New("Invalid or missing -properties-reader-uri flag")
		}
	}

	return nil
}
