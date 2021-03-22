package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

import (
	"context"
	"github.com/sfomuseum/go-flags/flagset"
	http_flags "github.com/whosonfirst/go-whosonfirst-spatial-http/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/server"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"log"
)

func main() {

	ctx := context.Background()

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = flags.AppendIndexingFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = http_flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	fs.Set("spatial-database-uri", "rtree://")
	fs.Set("properties-reader-uri", "mock://")

	flagset.Parse(fs)

	err = flags.ValidateCommonFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = flags.ValidateIndexingFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	err = http_flags.ValidateWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}
}
