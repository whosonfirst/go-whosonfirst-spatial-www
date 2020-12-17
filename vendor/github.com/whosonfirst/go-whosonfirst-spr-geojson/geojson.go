package geojson

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	go_geojson "github.com/paulmach/go.geojson"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"io"
)

func ToFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, r reader.Reader) (*go_geojson.FeatureCollection, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err := AsFeatureCollection(ctx, rsp, r, wr)

	if err != nil {
		return nil, err
	}

	wr.Flush()

	return go_geojson.UnmarshalFeatureCollection(buf.Bytes())
}

func ToFeatureCollectionWithJSON(ctx context.Context, body []byte, path string, r reader.Reader) (*go_geojson.FeatureCollection, error) {

	var buf bytes.Buffer
	wr := bufio.NewWriter(&buf)

	err := AsFeatureCollectionWithJSON(ctx, body, path, r, wr)

	if err != nil {
		return nil, err
	}

	wr.Flush()

	return go_geojson.UnmarshalFeatureCollection(buf.Bytes())
}

func AsFeatureCollection(ctx context.Context, rsp spr.StandardPlacesResults, r reader.Reader, wr io.Writer) error {

	fc, err := NewFeatureCollectionWriter(r, wr)

	if err != nil {
		return err
	}

	err = fc.Begin()

	if err != nil {
		return err
	}

	for _, pl := range rsp.Results() {

		err = fc.WriteFeature(ctx, pl.Path())

		if err != nil {
			return err
		}
	}

	return fc.End()
}

func AsFeatureCollectionWithJSON(ctx context.Context, body []byte, path string, r reader.Reader, wr io.Writer) error {

	path_rsp := gjson.GetBytes(body, path)

	if !path_rsp.Exists() {
		return errors.New("Missing path")
	}

	fc, err := NewFeatureCollectionWriter(r, wr)

	if err != nil {
		return err
	}

	err = fc.Begin()

	if err != nil {
		return err
	}

	for _, pl := range path_rsp.Array() {

		err := fc.WriteFeature(ctx, pl.String())

		if err != nil {
			return err
		}
	}

	return fc.End()
}
