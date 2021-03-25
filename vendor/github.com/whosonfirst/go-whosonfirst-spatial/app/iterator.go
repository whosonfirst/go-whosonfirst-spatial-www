package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/warning"
	"io"
	"log"
)

func NewIteratorWithFlagSet(ctx context.Context, fl *flag.FlagSet, spatial_db database.SpatialDatabase) (*iterator.Iterator, error) {

	emitter_uri, _ := lookup.StringVar(fl, flags.ITERATOR_URI)
	is_wof, _ := lookup.BoolVar(fl, flags.IS_WOF)

	emitter_cb := func(ctx context.Context, fh io.ReadSeeker, args ...interface{}) error {

		f, err := feature.LoadFeatureFromReader(fh)

		if err != nil {
			return err
		}

		if is_wof {

			if err != nil {

				// it's still not clear (to me) what the expected or desired
				// behaviour is / in this instance we might be issuing a warning
				// from the geojson-v2 package because a feature might have a
				// placetype defined outside of "core" (in the go-whosonfirst-placetypes)
				// package but that shouldn't necessarily trigger a fatal error
				// (20180405/thisisaaronland)

				if !warning.IsWarning(err) {
					return err
				}

				log.Printf("Feature ID %s triggered the following warning: %s\n", f.Id(), err)
			}
		}

		geom_type := geometry.Type(f)

		if geom_type == "Point" {
			return nil
		}

		err = spatial_db.IndexFeature(ctx, f)

		if err != nil {

			// something something something wrapping errors in Go 1.13
			// something something something waiting to see if the GOPROXY is
			// disabled by default in Go > 1.13 (20190919/thisisaaronland)

			msg := fmt.Sprintf("Failed to index %s (%s), %s", f.Id(), f.Name(), err)
			return errors.New(msg)
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, emitter_uri, emitter_cb)
	return iter, err
}
