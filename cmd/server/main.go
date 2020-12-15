package main

import (
	"context"
	_ "github.com/whosonfirst/go-whosonfirst-index/fs"
	http_flags "github.com/whosonfirst/go-whosonfirst-spatial-http/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/server"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-mock"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"log"
)

func main() {

	ctx := context.Background()

	fs, err := flags.CommonFlags()

	if err != nil {
		log.Fatal(err)
	}

	err = http_flags.AppendWWWFlags(fs)

	if err != nil {
		log.Fatal(err)
	}

	fs.Set("database-uri", "mock://")
	fs.Set("properties-reader-uri", "mock://")

	flags.Parse(fs)

	app, err := server.NewHTTPServerApplication(ctx)

	err = app.RunWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(err)
	}

}
