package geojson

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"github.com/whosonfirst/go-whosonfirst-spr"
	"github.com/whosonfirst/go-whosonfirst-uri"
)

type SPRPathResolver func(context.Context, spr.StandardPlacesResult) (string, error)

type JSONPathResolver func(context.Context, []byte) ([]string, error)

type JSONPathResolverCallback func(context.Context, string) (string, error)

func WhosOnFirstPathWithString(path string) (string, error) {

	id, uri_args, err := uri.ParseURI(path)

	if err != nil {
		return "", err
	}

	return uri.Id2RelPath(id, uri_args)
}

func WhosOnFirstSPRPathResolverFunc() SPRPathResolver {

	fn := func(ctx context.Context, r spr.StandardPlacesResult) (string, error) {
		return WhosOnFirstPathWithString(r.Id())
	}

	return fn
}

func JSONPathResolverFunc(gjson_path string) JSONPathResolver {

	return JSONPathResolverFuncWithCallback(gjson_path, nil)
}

func JSONPathResolverFuncWithCallback(gjson_path string, cb JSONPathResolverCallback) JSONPathResolver {

	fn := func(ctx context.Context, body []byte) ([]string, error) {

		path_rsp := gjson.GetBytes(body, gjson_path)

		if !path_rsp.Exists() {
			return nil, errors.New("Missing path")
		}

		paths := make([]string, 0)

		for _, p := range path_rsp.Array() {

			path := p.String()

			if cb != nil {

				new_path, err := cb(ctx, path)

				if err != nil {
					return nil, err
				}

				path = new_path
			}

			paths = append(paths, path)
		}

		return paths, nil
	}

	return fn
}
