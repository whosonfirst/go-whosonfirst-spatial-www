package properties

import (
	"context"
	"github.com/aaronland/go-roster"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net/url"
)

type PropertiesResponse map[string]interface{}

type PropertiesResponseResults struct {
	Properties []*PropertiesResponse `json:"properties"`
}

type PropertiesReader interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PropertiesResponseResultsWithStandardPlacesResults(context.Context, spr.StandardPlacesResults, []string) (*PropertiesResponseResults, error)
	AppendPropertiesWithFeatureCollection(context.Context, *geojson.GeoJSONFeatureCollection, []string) error
	Close(context.Context) error
}

type PropertiesReaderInitializeFunc func(ctx context.Context, uri string) (PropertiesReader, error)

var properties_readers roster.Roster

func ensurePropertiesRoster() error {

	if properties_readers == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		properties_readers = r
	}

	return nil
}

func RegisterPropertiesReader(ctx context.Context, scheme string, f PropertiesReaderInitializeFunc) error {

	err := ensurePropertiesRoster()

	if err != nil {
		return err
	}

	return properties_readers.Register(ctx, scheme, f)
}

func NewPropertiesReader(ctx context.Context, uri string) (PropertiesReader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := properties_readers.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(PropertiesReaderInitializeFunc)
	return f(ctx, uri)
}
