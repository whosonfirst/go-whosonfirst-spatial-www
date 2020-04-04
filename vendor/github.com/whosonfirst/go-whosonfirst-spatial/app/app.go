package app

import (
	"context"
	"flag"
	"github.com/whosonfirst/go-whosonfirst-index"
	"github.com/whosonfirst/go-whosonfirst-log"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"runtime/debug"
	"time"
)

type SpatialApplication struct {
	mode             string
	SpatialDatabase  database.SpatialDatabase
	PropertiesReader properties.PropertiesReader
	Walker           *index.Indexer
	Logger           *log.WOFLogger
}

func NewSpatialApplicationWithFlagSet(ctx context.Context, fl *flag.FlagSet) (*SpatialApplication, error) {

	logger, err := NewApplicationLoggerWithFlagSet(ctx, fl)

	if err != nil {
		return nil, err
	}

	spatial_db, err := NewSpatialDatabaseWithFlagSet(ctx, fl)

	if err != nil {
		return nil, err
	}

	properties_r, err := NewPropertiesReaderWithFlagSet(ctx, fl)

	if err != nil {
		return nil, err
	}

	walker, err := NewWalkerWithFlagSet(ctx, fl, spatial_db, properties_r)

	if err != nil {
		return nil, err
	}

	err = AppendCustomPlacetypesWithFlagSet(ctx, fl)

	if err != nil {
		return nil, err
	}

	sp := SpatialApplication{
		SpatialDatabase:  spatial_db,
		PropertiesReader: properties_r,
		Walker:           walker,
		Logger:           logger,
	}

	return &sp, nil
}

func (p *SpatialApplication) Close(ctx context.Context) error {

	p.SpatialDatabase.Close(ctx)

	if p.PropertiesReader != nil {
		p.PropertiesReader.Close(ctx)
	}

	return nil
}

func (p *SpatialApplication) IndexPaths(ctx context.Context, paths ...string) error {

	if p.mode != "spatialite" {

		go func() {

			// TO DO: put this somewhere so that it can be triggered by signal(s)
			// to reindex everything in bulk or incrementally

			t1 := time.Now()

			err := p.Walker.IndexPaths(paths)

			if err != nil {
				p.Logger.Fatal("failed to index paths because %s", err)
			}

			t2 := time.Since(t1)

			p.Logger.Status("finished indexing in %v", t2)
			debug.FreeOSMemory()
		}()

		// set up some basic monitoring and feedback stuff

		go func() {

			c := time.Tick(1 * time.Second)

			for _ = range c {

				if !p.Walker.IsIndexing() {
					continue
				}

				p.Logger.Status("indexing %d records indexed", p.Walker.Indexed)
			}
		}()
	}

	return nil
}

/*

	go func() {

		tick := time.Tick(1 * time.Minute)

		for _ = range tick {
			var ms runtime.MemStats
			runtime.ReadMemStats(&ms)
			pip.Logger.Status("memstats system: %8d inuse: %8d released: %8d objects: %6d", ms.HeapSys, ms.HeapInuse, ms.HeapReleased, ms.HeapObjects)
		}
	}()

*/
