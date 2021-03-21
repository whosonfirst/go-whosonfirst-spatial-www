package app

import (
	"context"
	"flag"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"io"
	"os"
)

func NewApplicationLoggerWithFlagSet(ctx context.Context, fl *flag.FlagSet) (*log.WOFLogger, error) {

	verbose, _ := lookup.BoolVar(fl, flags.VERBOSE)

	logger := log.SimpleWOFLogger()
	level := "status"

	if verbose {
		level = "debug"
	}

	stdout := io.Writer(os.Stdout)
	logger.AddLogger(stdout, level)

	return logger, nil
}
