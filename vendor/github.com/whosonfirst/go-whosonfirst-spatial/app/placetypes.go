package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"io"
	"strings"
)

func AppendCustomPlacetypesWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	enable_custom_placetypes, _ := flags.BoolVar(fs, "enable-custom-placetypes")
	custom_placetypes_source, _ := flags.StringVar(fs, "custom-placetypes-source")
	custom_placetypes, _ := flags.StringVar(fs, "custom-placetypes")

	if !enable_custom_placetypes {
		return nil
	}

	var custom_reader io.Reader

	if custom_placetypes_source == "" {
		custom_reader = strings.NewReader(custom_placetypes)
	} else {

		// whosonfirst/go-reader or ... ?
	}

	spec, err := placetypes.NewWOFPlacetypeSpecificationWithReader(custom_reader)

	if err != nil {
		return err
	}

	return placetypes.AppendPlacetypeSpecification(spec)
}
