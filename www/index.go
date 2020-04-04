package www

import (
	"errors"
	"html/template"
	gohttp "net/http"
)

type IndexHandlerOptions struct {
	Templates *template.Template
}

func IndexHandler(opts *IndexHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("index")

	if t == nil {
		return nil, errors.New("Missing 'index' template")
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}
