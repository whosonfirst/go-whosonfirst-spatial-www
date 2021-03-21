package app

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
)

func NewPropertiesReaderWithFlagSet(ctx context.Context, fl *flag.FlagSet) (properties.PropertiesReader, error) {

	enable_properties, _ := lookup.BoolVar(fl, flags.ENABLE_PROPERTIES)
	properties_reader_uri, _ := lookup.StringVar(fl, flags.PROPERTIES_READER_URI)

	if !enable_properties {
		return nil, nil
	}

	return properties.NewPropertiesReader(ctx, properties_reader_uri)
}
