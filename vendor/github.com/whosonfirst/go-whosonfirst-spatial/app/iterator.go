package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/geometry"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
	"github.com/whosonfirst/go-whosonfirst-iterate/iterator"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/warning"
	"io"
	"log"
	"strings"
	"sync"
)

func NewIteratorWithFlagSet(ctx context.Context, fl *flag.FlagSet, spatial_db database.SpatialDatabase, properties_r properties.PropertiesReader) (*iterator.Iterator, error) {

	emitter_uri, _ := lookup.StringVar(fl, flags.ITERATOR_URI)

	is_wof, _ := lookup.BoolVar(fl, "is-wof")
	index_properties, _ := lookup.BoolVar(fl, flags.INDEX_PROPERTIES)

	include_deprecated := true
	include_superseded := true
	include_ceased := true
	include_notcurrent := true

	exclude_fl := fl.Lookup("exclude")

	if exclude_fl != nil {

		// ugh... Go - why do I have to do this... I am willing
		// to believe I am "doing it wrong" (obviously) but for
		// the life of me I can't figure out how to do it "right"
		// (20180301/thisisaaronland)

		exclude := strings.Split(exclude_fl.Value.String(), " ")

		for _, e := range exclude {

			switch e {
			case "deprecated":
				include_deprecated = false
			case "ceased":
				include_ceased = false
			case "superseded":
				include_superseded = false
			case "not-current":
				include_notcurrent = false
			default:
				continue
			}
		}
	}

	var wg *sync.WaitGroup
	var mu *sync.Mutex

	if index_properties {
		wg = new(sync.WaitGroup)
		mu = new(sync.Mutex)
	}

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

			if !include_notcurrent {

				fl, err := whosonfirst.IsCurrent(f)

				if err != nil {
					return err
				}

				if fl.IsTrue() && fl.IsKnown() {
					return nil
				}
			}

			if !include_deprecated {

				fl, err := whosonfirst.IsDeprecated(f)

				if err != nil {
					return err
				}

				if fl.IsTrue() && fl.IsKnown() {
					return nil
				}
			}

			if !include_ceased {

				fl, err := whosonfirst.IsCeased(f)

				if err != nil {
					return err
				}

				if fl.IsTrue() && fl.IsKnown() {
					return nil
				}
			}

			if !include_superseded {

				fl, err := whosonfirst.IsSuperseded(f)

				if err != nil {
					return err
				}

				if fl.IsTrue() && fl.IsKnown() {
					return nil
				}
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

		if index_properties {

			wg.Add(1)

			go func(f geojson.Feature, wg *sync.WaitGroup) error {

				defer wg.Done()

				mu.Lock()

				err = properties_r.IndexFeature(ctx, f)

				mu.Unlock()

				if err != nil {
					// log.Println("FAILED TO INDEX", err) // something
				}

				return err

			}(f, wg)
		}

		return nil
	}

	iter, err := iterator.NewIterator(ctx, emitter_uri, emitter_cb)

	if index_properties {
		wg.Wait()
	}

	return iter, err
}
