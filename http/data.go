package http

// TBD: make this part of whosonfirst/go-reader package...

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	gohttp "net/http"
)

func NewDataHandler(r reader.Reader) (gohttp.Handler, error) {

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		path := req.URL.Path

		id, uri_args, err := uri.ParseURI(path)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		ctx := req.Context()
		fh, err := r.Read(ctx, rel_path)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		_, err = io.Copy(rsp, fh)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusBadRequest)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
