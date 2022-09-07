package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
	_ "github.com/whosonfirst/go-reader-cachereader"	
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/application/server"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)
	
	if err != nil {
		logger.Fatal(err)
	}
}
