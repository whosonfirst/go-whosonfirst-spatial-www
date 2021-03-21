package spatial

import (
	"context"
	"github.com/skelterjohn/geom"
	"github.com/whosonfirst/go-whosonfirst-flags"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type SpatialIndex interface {
	IndexFeature(context.Context, wof_geojson.Feature) error
	PointInPolygon(context.Context, *geom.Coord, ...Filter) (spr.StandardPlacesResults, error)
	PointInPolygonCandidates(context.Context, *geom.Coord, ...Filter) ([]*PointInPolygonCandidate, error)
	PointInPolygonWithChannels(context.Context, chan spr.StandardPlacesResult, chan error, chan bool, *geom.Coord, ...Filter)
	PointInPolygonCandidatesWithChannels(context.Context, chan *PointInPolygonCandidate, chan error, chan bool, *geom.Coord, ...Filter)
	Disconnect(context.Context) error
}

type PointInPolygonCandidate struct {
	Id        string
	FeatureId string
	IsAlt     bool
	AltLabel  string
	Bounds    *geom.Rect
}

type PropertiesResponse map[string]interface{}

type PropertiesResponseResults struct {
	Properties []*PropertiesResponse `json:"properties"`
}

type Filter interface {
	HasPlacetypes(flags.PlacetypeFlag) bool
	MatchesInception(flags.DateFlag) bool
	MatchesCessation(flags.DateFlag) bool
	IsCurrent(flags.ExistentialFlag) bool
	IsDeprecated(flags.ExistentialFlag) bool
	IsCeased(flags.ExistentialFlag) bool
	IsSuperseded(flags.ExistentialFlag) bool
	IsSuperseding(flags.ExistentialFlag) bool
	IsAlternateGeometry(flags.AlternateGeometryFlag) bool
	HasAlternateGeometry(flags.AlternateGeometryFlag) bool
}
