package data

import (
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-uri"
	"io"
	"net/http"
)

type WhosOnFirstDataHandlerOptions struct {
	ContentType string
}

func DefaultWhosOnFirstDataHandlerOptions() *WhosOnFirstDataHandlerOptions {

	opts := &WhosOnFirstDataHandlerOptions{
		ContentType: "application/json",
	}

	return opts
}

func WhosOnFirstDataHandler(r reader.Reader) http.Handler {

	opts := DefaultWhosOnFirstDataHandlerOptions()
	return WhosOnFirstDataHandlerWithOptions(r, opts)
}

func WhosOnFirstDataHandlerWithOptions(r reader.Reader, opts *WhosOnFirstDataHandlerOptions) http.Handler {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		path := req.URL.Path

		id, uri_args, err := uri.ParseURI(path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		rel_path, err := uri.Id2RelPath(id, uri_args)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := req.Context()

		fh, err := r.Read(ctx, rel_path)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		defer fh.Close()

		rsp.Header().Set("Content-type", opts.ContentType)
		
		_, err = io.Copy(rsp, fh)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	return http.HandlerFunc(fn)
}
