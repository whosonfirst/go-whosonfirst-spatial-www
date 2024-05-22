package main

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
)

import (
	"context"
	"log/slog"
	"os"

	"github.com/whosonfirst/go-whosonfirst-spatial-www/app/server"
)

func main() {

	ctx := context.Background()
	logger := slog.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
