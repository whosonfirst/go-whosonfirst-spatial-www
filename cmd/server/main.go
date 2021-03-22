package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

import (
	"context"
	"github.com/sfomuseum/go-flags/flagset"
	www_flags "github.com/whosonfirst/go-whosonfirst-spatial-www/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/server"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"log"
)

func main() {

	ctx := context.Background()

	fs, err := spatial_flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = spatial_flags.AppendIndexingFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = www_flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	flagset.Parse(fs)

	err = spatial_flags.ValidateCommonFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = spatial_flags.ValidateIndexingFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = www_flags.ValidateWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}
}
