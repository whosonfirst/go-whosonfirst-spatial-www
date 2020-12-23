package spatial

import (
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/properties/whosonfirst"
)

func SpatialIdWithFeature(f geojson.Feature, extra ...interface{}) (string, error) {

	feature_id := f.Id()
	alt_label := whosonfirst.AltLabel(f)

	sp_id := fmt.Sprintf("%s#%s", feature_id, alt_label)

	if len(extra) > 0 {

		for _, v := range extra {
			sp_id = fmt.Sprintf("%s:%v", sp_id, v)
		}
	}
	return sp_id, nil
}
