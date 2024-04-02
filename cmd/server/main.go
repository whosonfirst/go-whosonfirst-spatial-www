package main

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

import (
	"context"
	"log"

	"github.com/whosonfirst/go-whosonfirst-spatial-www/app/server"	
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatal(err)
	}
}
