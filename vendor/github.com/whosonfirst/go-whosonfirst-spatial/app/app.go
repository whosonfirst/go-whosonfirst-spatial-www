package app

import (
	"context"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-iterate/v2/iterator"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"log"
	"runtime/debug"
	"time"
)

type SpatialApplication struct {
	mode             string
	SpatialDatabase  database.SpatialDatabase
	PropertiesReader reader.Reader
	Iterator         *iterator.Iterator
	Logger           *log.Logger
}

func NewSpatialApplicationWithFlagSet(ctx context.Context, fl *flag.FlagSet) (*SpatialApplication, error) {

	logger, err := NewApplicationLoggerWithFlagSet(ctx, fl)

	if err != nil {
		return nil, err
	}

	spatial_db, err := NewSpatialDatabaseWithFlagSet(ctx, fl)

	if err != nil {
		return nil, fmt.Errorf("Failed instantiate spatial database, %v", err)
	}

	properties_r, err := NewPropertiesReaderWithFlagsSet(ctx, fl)

	if err != nil {
		return nil, fmt.Errorf("Failed to create properties reader, %v", err)
	}

	if properties_r == nil {
		properties_r = spatial_db
	}

	iter, err := NewIteratorWithFlagSet(ctx, fl, spatial_db)

	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate iterator, %v", err)
	}

	err = AppendCustomPlacetypesWithFlagSet(ctx, fl)

	if err != nil {
		return nil, fmt.Errorf("Failed to append custom placetypes, %v", err)
	}

	sp := SpatialApplication{
		SpatialDatabase:  spatial_db,
		PropertiesReader: properties_r,
		Iterator:         iter,
		Logger:           logger,
	}

	return &sp, nil
}

func (p *SpatialApplication) Close(ctx context.Context) error {

	p.SpatialDatabase.Disconnect(ctx)

	return nil
}

func (p *SpatialApplication) IndexPaths(ctx context.Context, paths ...string) error {

	go func() {

		// TO DO: put this somewhere so that it can be triggered by signal(s)
		// to reindex everything in bulk or incrementally

		t1 := time.Now()

		err := p.Iterator.IterateURIs(ctx, paths...)

		if err != nil {
			p.Logger.Fatalf("failed to index paths because %s", err)
		}

		t2 := time.Since(t1)

		p.Logger.Printf("finished indexing in %v", t2)
		debug.FreeOSMemory()
	}()

	// set up some basic monitoring and feedback stuff

	go func() {

		c := time.Tick(1 * time.Second)

		for _ = range c {

			if !p.Iterator.IsIndexing() {
				continue
			}

			p.Logger.Printf("indexing %d records indexed", p.Iterator.Seen)
		}
	}()

	return nil
}
