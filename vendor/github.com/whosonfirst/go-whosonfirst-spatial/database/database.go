package database

import (
	"context"
	"github.com/aaronland/go-roster"
	"github.com/skelterjohn/geom"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"net/url"
)

type SpatialDatabase interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PointInPolygon(context.Context, *geom.Coord, filter.Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, *geom.Coord) (*geojson.GeoJSONFeatureCollection, error)
	PointInPolygonWithChannels(context.Context, *geom.Coord, filter.Filter, chan spr.StandardPlacesResult, chan error, chan bool)
	PointInPolygonCandidatesWithChannels(context.Context, *geom.Coord, chan geojson.GeoJSONFeature, chan error, chan bool)
	StandardPlacesResultsToFeatureCollection(context.Context, spr.StandardPlacesResults) (*geojson.GeoJSONFeatureCollection, error)
	Close(context.Context) error
}

type SpatialDatabaseInitializeFunc func(ctx context.Context, uri string) (SpatialDatabase, error)

var spatial_databases roster.Roster

func ensureSpatialRoster() error {

	if spatial_databases == nil {

		r, err := roster.NewDefaultRoster()

		if err != nil {
			return err
		}

		spatial_databases = r
	}

	return nil
}

func RegisterSpatialDatabase(ctx context.Context, scheme string, f SpatialDatabaseInitializeFunc) error {

	err := ensureSpatialRoster()

	if err != nil {
		return err
	}

	return spatial_databases.Register(ctx, scheme, f)
}

func NewSpatialDatabase(ctx context.Context, uri string) (SpatialDatabase, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	scheme := u.Scheme

	i, err := spatial_databases.Driver(ctx, scheme)

	if err != nil {
		return nil, err
	}

	f := i.(SpatialDatabaseInitializeFunc)
	return f(ctx, uri)
}
