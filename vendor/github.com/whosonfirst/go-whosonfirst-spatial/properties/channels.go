package properties

import (
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
)

type ChannelResponse struct {
	Index   int
	Feature geojson.GeoJSONFeature
}
