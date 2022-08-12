package main

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/server"
	"log"
)

func main() {

	ctx := context.Background()

	app, err := server.NewHTTPServerApplication(ctx)

	if err != nil {
		log.Fatal(err)
	}

	err = app.Run(ctx)

	if err != nil {
		log.Fatal(err)
	}
}
