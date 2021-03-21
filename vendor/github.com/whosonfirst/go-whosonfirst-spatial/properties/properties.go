package properties

import (
	"context"
	"fmt"
	"github.com/aaronland/go-roster"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"net/url"
	"sort"
	"strings"
)

type PropertiesReader interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PropertiesResponseResultsWithStandardPlacesResults(context.Context, spr.StandardPlacesResults, []string) (*spatial.PropertiesResponseResults, error)
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

func Schemes() []string {

	ctx := context.Background()
	schemes := []string{}

	err := ensurePropertiesRoster()

	if err != nil {
		return schemes
	}

	for _, dr := range properties_readers.Drivers(ctx) {
		scheme := fmt.Sprintf("%s://", strings.ToLower(dr))
		schemes = append(schemes, scheme)
	}

	sort.Strings(schemes)
	return schemes
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
