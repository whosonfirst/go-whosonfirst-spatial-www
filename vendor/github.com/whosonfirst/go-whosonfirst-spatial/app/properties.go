package app

import (
	_ "github.com/whosonfirst/go-reader-http"
	_ "github.com/whosonfirst/go-reader-whosonfirst-data"	
)

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
)

func NewPropertiesReaderWithFlagSet(ctx context.Context, fl *flag.FlagSet) (properties.PropertiesReader, error) {

	enable_properties, _ := flags.BoolVar(fl, "enable-properties")
	properties_reader_uri, _ := flags.StringVar(fl, "properties-reader-uri")

	if !enable_properties {
		return nil, nil
	}

	return properties.NewPropertiesReader(ctx, properties_reader_uri)
}
