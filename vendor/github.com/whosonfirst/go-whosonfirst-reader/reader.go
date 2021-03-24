package reader

import (
	"context"
	go_reader "github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2"
	"github.com/whosonfirst/go-whosonfirst-geojson-v2/feature"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"io/ioutil"
)

func LoadReadCloserFromID(ctx context.Context, r go_reader.Reader, id int64) (io.ReadCloser, error) {

	rel_path, err := uri.Id2RelPath(id)

	if err != nil {
		return nil, err
	}

	return r.Read(ctx, rel_path)
}

func LoadBytesFromID(ctx context.Context, r go_reader.Reader, id int64) ([]byte, error) {

	fh, err := LoadReadCloserFromID(ctx, r, id)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return ioutil.ReadAll(fh)
}

func LoadFeatureFromID(ctx context.Context, r go_reader.Reader, id int64) (geojson.Feature, error) {

	fh, err := LoadReadCloserFromID(ctx, r, id)

	if err != nil {
		return nil, err
	}

	defer fh.Close()

	return feature.LoadWOFFeatureFromReader(fh)
}
