package properties

import (
	"github.com/paulmach/go.geojson"
)

type ChannelResponse struct {
	Index   int
	Feature *geojson.Feature
}
