package mock

import (
	"context"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	spatial_properties "github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/paulmach/go.geojson"	
)

type MockPropertiesReader struct {
	spatial_properties.PropertiesReader
}

func init() {
	ctx := context.Background()
	spatial_properties.RegisterPropertiesReader(ctx, "mock", NewMockPropertiesReader)
}

func NewMockPropertiesReader(ctx context.Context, uri string) (spatial_properties.PropertiesReader, error) {
	pr := &MockPropertiesReader{}
	return pr, nil
}

func (pr *MockPropertiesReader) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {
	return nil
}

func (pr *MockPropertiesReader) PropertiesResponseResultsWithStandardPlacesResults(ctx context.Context, results spr.StandardPlacesResults, properties []string) (*spatial_properties.PropertiesResponseResults, error) {
	return nil, nil
}

func (pr *MockPropertiesReader) AppendPropertiesWithFeatureCollection(ctx context.Context, fc *geojson.FeatureCollection, properties []string) error {
	return nil
}

func (pr *MockPropertiesReader) Close(ctx context.Context) error {
	return nil
}
