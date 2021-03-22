package rtree

// Really this is a kind of 'null' properties reader but it
// might be updated at some later date to have an in-memory
// cache of a finite set of properties. TBD...
// (20210322/straup)

import (
	"context"
	"encoding/json"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	spatial_properties "github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
)

type RTreePropertiesReader struct {
	spatial_properties.PropertiesReader
}

func init() {
	ctx := context.Background()
	spatial_properties.RegisterPropertiesReader(ctx, "rtree", NewRTreePropertiesReader)
}

func NewRTreePropertiesReader(ctx context.Context, uri string) (spatial_properties.PropertiesReader, error) {

	pr := &RTreePropertiesReader{}

	return pr, nil
}

func (pr *RTreePropertiesReader) IndexFeature(ctx context.Context, f wof_geojson.Feature) error {
	// See notes above about potentially adding a fixed-list of
	// properties to store in-memory (20210322/straup)
	return nil
}

func (pr *RTreePropertiesReader) PropertiesResponseResultsWithStandardPlacesResults(ctx context.Context, results spr.StandardPlacesResults, properties []string) (*spatial.PropertiesResponseResults, error) {

	previous_results := results.Results()

	new_results := make([]*spatial.PropertiesResponse, len(previous_results))

	for idx, r := range previous_results {

		target, err := json.Marshal(r)

		source := []byte(``)

		append_opts := &spatial_properties.AppendPropertiesOptions{
			Keys:         properties,
			SourcePrefix: "",
			TargetPrefix: "",
		}

		target, err = spatial_properties.AppendPropertiesWithJSON(ctx, append_opts, source, target)

		if err != nil {
			return nil, err
		}

		var props *spatial.PropertiesResponse
		err = json.Unmarshal(target, &props)

		if err != nil {
			return nil, err
		}

		new_results[idx] = props
	}

	props_rsp := &spatial.PropertiesResponseResults{
		Properties: new_results,
	}

	return props_rsp, nil
}

func (pr *RTreePropertiesReader) Close(ctx context.Context) error {
	return nil
}
