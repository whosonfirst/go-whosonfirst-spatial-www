package flags

import (
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
)

func CommonFlags() (*flag.FlagSet, error) {

	fs := flagset.NewFlagSet("common")

	// spatial databases

	available_databases := database.Schemes()
	desc_databases := fmt.Sprintf("Valid options are: %s", available_databases)

	fs.String(SPATIAL_DATABASE_URI, "", desc_databases)

	fs.Bool(IS_WOF, true, "Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents.")

	fs.Bool(ENABLE_CUSTOM_PLACETYPES, false, "Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.")

	// Pending changes in the app/placetypes.go package to support
	// alternate sources (20210324/thisisaaronland)
	// fs.String(CUSTOM_PLACETYPES_SOURCE, "", "...")

	fs.String(CUSTOM_PLACETYPES, "", "A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.")

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

	return nil
}
