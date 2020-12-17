package mock

import (
	"context"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

type MockPropertiesReader struct {
	properties.PropertiesReader
}

func init() {
	ctx := context.Background()
	properties.RegisterPropertiesReader(ctx, "mock", NewMockPropertiesReader)
}

func NewMockPropertiesReader(ctx context.Context, uri string) (properties.PropertiesReader, error) {
	pr := &MockPropertiesReader{}
	return pr, nil
}

func (pr *MockPropertiesReader) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {
	return nil
}

func (pr *MockPropertiesReader) PropertiesResponseResultsWithStandardPlacesResults(ctx context.Context, results spr.StandardPlacesResults, properties []string) (*spatial.PropertiesResponseResults, error) {
	return nil, nil
}

func (pr *MockPropertiesReader) Close(ctx context.Context) error {
	return nil
}
