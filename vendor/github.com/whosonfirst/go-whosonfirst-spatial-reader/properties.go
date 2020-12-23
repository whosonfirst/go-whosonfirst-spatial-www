package reader

import (
	"context"
	"encoding/json"
	"errors"
	go_reader "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial/properties"
	"github.com/whosonfirst/go-whosonfirst-spr"
	_ "log"
	"net/url"
	"strconv"
)

func init() {
	ctx := context.Background()
	properties.RegisterPropertiesReader(ctx, "whosonfirst", NewWhosonfirstPropertiesReader)
}

type WhosonfirstPropertiesReader struct {
	properties.PropertiesReader
	reader go_reader.Reader
}

func NewWhosonfirstPropertiesReader(ctx context.Context, uri string) (properties.PropertiesReader, error) {

	u, err := url.Parse(uri)

	if err != nil {
		return nil, err
	}

	q := u.Query()

	reader_uri := q.Get("reader")

	if reader_uri == "" {
		return nil, errors.New("Missing reader parameter")
	}

	cache_uri := q.Get("cache")

	if cache_uri == "" {
		cache_uri = "null://"
	}

	cr_q := url.Values{}
	cr_q.Set("reader", reader_uri)
	cr_q.Set("cache", cache_uri)

	cr_uri := url.URL{}
	cr_uri.Scheme = "cachereader"
	cr_uri.RawQuery = cr_q.Encode()

	cr, err := cachereader.NewCacheReader(ctx, cr_uri.String())

	if err != nil {
		return nil, err
	}

	db := &WhosonfirstPropertiesReader{
		reader: cr,
	}

	return db, nil
}

func (db *WhosonfirstPropertiesReader) Close(ctx context.Context) error {
	return nil
}

func (db *WhosonfirstPropertiesReader) IndexFeature(context.Context, wof_geojson.Feature) error {
	return nil
}

func (db *WhosonfirstPropertiesReader) PropertiesResponseResultsWithStandardPlacesResults(ctx context.Context, results spr.StandardPlacesResults, property_keys []string) (*spatial.PropertiesResponseResults, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	previous_results := results.Results()
	new_results := make([]*spatial.PropertiesResponse, len(previous_results))

	for idx, r := range previous_results {

		target, err := json.Marshal(r)

		if err != nil {
			return nil, err
		}

		str_id := r.Id()

		id, err := strconv.ParseInt(str_id, 10, 64)

		if err != nil {
			return nil, err
		}

		source, err := wof_reader.LoadBytesFromID(ctx, db.reader, id)

		if err != nil {
			return nil, err
		}

		append_opts := &properties.AppendPropertiesOptions{
			Keys:         property_keys,
			SourcePrefix: "properties",
			TargetPrefix: "",
		}

		target, err = properties.AppendPropertiesWithJSON(ctx, append_opts, source, target)

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
