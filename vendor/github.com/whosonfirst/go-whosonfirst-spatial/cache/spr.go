package cache

import (
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/geometry"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
)

// see the way we're storing a geojson.Geometry but returning a
// geojson.Polygons per the interface definition? see notes in
// go-whosonfirst-geojson-v2/geometry/polygon.go function called
// PolygonsForFeature for details (20170921/thisisaaronland)

type SPRCacheItem struct {
	CacheItem       `json:",omitempty"`
	FeatureSPR      spr.StandardPlacesResult `json:"spr"`
	FeaturePolygons []wof_geojson.Polygon    `json:"polygons"`
	// FeatureSPR      *feature.WOFStandardPlacesResult  `json:"spr"`
	// FeaturePolygons []geometry.Polygon    `json:"polygons"`
}

func NewSPRCacheItem(f wof_geojson.Feature) (CacheItem, error) {

	s, err := f.SPR()

	if err != nil {
		return nil, err
	}

	polys, err := geometry.PolygonsForFeature(f)

	if err != nil {
		return nil, err
	}

	fc := SPRCacheItem{
		// FeatureSPR:      s.(*feature.WOFStandardPlacesResult),
		FeatureSPR:      s,
		FeaturePolygons: polys,
	}

	return &fc, nil
}

func (fc *SPRCacheItem) SPR() spr.StandardPlacesResult {
	return fc.FeatureSPR
}

func (fc *SPRCacheItem) Geometry() geojson.GeoJSONGeometry {

	multi_poly := make([]geojson.GeoJSONPolygon, 0)

	for _, p := range fc.Polygons() {

		poly := make([]geojson.GeoJSONRing, 0)

		ext := p.ExteriorRing()

		ext_ring := make([]geojson.GeoJSONPoint, 0)

		for _, coord := range ext.Vertices() {

			pt := geojson.GeoJSONPoint{coord.X, coord.Y}
			ext_ring = append(ext_ring, pt)
		}

		poly = append(poly, ext_ring)

		for _, int := range p.InteriorRings() {

			int_ring := make([]geojson.GeoJSONPoint, 0)

			for _, coord := range int.Vertices() {

				pt := geojson.GeoJSONPoint{coord.X, coord.Y}
				int_ring = append(int_ring, pt)
			}

			poly = append(poly, int_ring)
		}

		multi_poly = append(multi_poly, poly)
	}

	geom := geojson.GeoJSONGeometry{
		Type:        "MultiPolygon",
		Coordinates: multi_poly,
	}

	return geom
}

func (fc *SPRCacheItem) Polygons() []wof_geojson.Polygon {

	return fc.FeaturePolygons
}
