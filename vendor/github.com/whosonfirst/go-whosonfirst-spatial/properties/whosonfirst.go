package properties

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-cache"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-reader-cachereader"
	wof_geojson "github.com/whosonfirst/go-whosonfirst-geojson-v2"
	wof_reader "github.com/whosonfirst/go-whosonfirst-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial/geojson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	_ "log"
	"net/url"
	"strconv"
)

func init() {
	ctx := context.Background()
	RegisterPropertiesReader(ctx, "whosonfirst", NewWhosonfirstPropertiesReader)
}

type ChannelResponse struct {
	Index   int
	Feature geojson.GeoJSONFeature
}

type WhosonfirstPropertiesReader struct {
	PropertiesReader
	reader reader.Reader
}

func NewWhosonfirstPropertiesReader(ctx context.Context, uri string) (PropertiesReader, error) {

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

	r, err := reader.NewReader(ctx, reader_uri)

	if err != nil {
		return nil, err
	}

	c, err := cache.NewCache(ctx, cache_uri)

	if err != nil {
		return nil, err
	}

	cr, err := cachereader.NewCacheReader(r, c)

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

func (db *WhosonfirstPropertiesReader) PropertiesResponseResultsWithStandardPlacesResults(ctx context.Context, results spr.StandardPlacesResults, properties []string) (*PropertiesResponseResults, error) {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	previous_results := results.Results()
	new_results := make([]*PropertiesResponse, len(previous_results))

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

		target, err = AppendPropertiesWithJSON(ctx, source, target, properties, "")

		if err != nil {
			return nil, err
		}

		var props *PropertiesResponse
		err = json.Unmarshal(target, &props)

		if err != nil {
			return nil, err
		}

		new_results[idx] = props
	}

	props_rsp := &PropertiesResponseResults{
		Properties: new_results,
	}

	return props_rsp, nil
}

func (db *WhosonfirstPropertiesReader) AppendPropertiesWithFeatureCollection(ctx context.Context, fc *geojson.GeoJSONFeatureCollection, properties []string) error {

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	rsp_ch := make(chan ChannelResponse)
	err_ch := make(chan error)
	done_ch := make(chan bool)

	remaining := len(fc.Features)

	for idx, f := range fc.Features {
		go db.appendPropertiesWithChannels(ctx, idx, f, properties, rsp_ch, err_ch, done_ch)
	}

	for remaining > 0 {
		select {
		case <-ctx.Done():
			return nil
		case <-done_ch:
			remaining -= 1
		case rsp := <-rsp_ch:
			fc.Features[rsp.Index] = rsp.Feature
		case err := <-err_ch:
			return err
		default:
			// pass
		}
	}

	return nil
}

func (db *WhosonfirstPropertiesReader) appendPropertiesWithChannels(ctx context.Context, idx int, f geojson.GeoJSONFeature, properties []string, rsp_ch chan ChannelResponse, err_ch chan error, done_ch chan bool) {

	defer func() {
		done_ch <- true
	}()

	select {
	case <-ctx.Done():
		return
	default:
		// pass
	}

	target, err := json.Marshal(f)

	if err != nil {
		err_ch <- err
		return
	}

	id_rsp := gjson.GetBytes(target, "properties.wof:id")

	if !id_rsp.Exists() {
		err_ch <- errors.New("Missing wof:id")
		return
	}

	id := id_rsp.Int()

	source, err := wof_reader.LoadBytesFromID(ctx, db.reader, id)

	if err != nil {
		err_ch <- err
		return
	}

	target, err = AppendPropertiesWithJSON(ctx, source, target, properties, "properties")

	if err != nil {
		err_ch <- err
		return
	}

	var new_f geojson.GeoJSONFeature
	err = json.Unmarshal(target, &new_f)

	if err != nil {
		err_ch <- err
		return
	}

	rsp := ChannelResponse{
		Index:   idx,
		Feature: new_f,
	}

	rsp_ch <- rsp
	return
}
