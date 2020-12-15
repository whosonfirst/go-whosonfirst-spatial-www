package mock

import (
	"context"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/paulmach/go.geojson"	
)

func init() {
	ctx := context.Background()
	database.RegisterSpatialDatabase(ctx, "mock", NewMockSpatialDatabase)
}

type MockSpatialDatabase struct {
	database.SpatialDatabase
}

func NewMockSpatialDatabase(ctx context.Context, uri string) (database.SpatialDatabase, error){
	db := &MockSpatialDatabase{}
	return db, nil
}

func (db *MockSpatialDatabase) Close(ctx context.Context) error {
	return nil
}

func (db *MockSpatialDatabase) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {
	return nil
}

func (db *MockSpatialDatabase) PointInPolygon(ctx context.Context, coord *geom.Coord, filters ...filter.Filter) (spr.StandardPlacesResults, error) {
	return nil, nil
}

func (db *MockSpatialDatabase) PointInPolygonWithChannels(ctx context.Context, rsp_ch chan spr.StandardPlacesResult, err_ch chan error, done_ch chan bool, coord *geom.Coord, filters ...filter.Filter) {
	return
}

func (db *MockSpatialDatabase) PointInPolygonCandidates(ctx context.Context, coord *geom.Coord) (*geojson.FeatureCollection, error) {
	return nil, nil
}

func (db *MockSpatialDatabase) PointInPolygonCandidatesWithChannels(ctx context.Context, coord *geom.Coord, rsp_ch chan *geojson.Feature, err_ch chan error, done_ch chan bool) {
	return
}
