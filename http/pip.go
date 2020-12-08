package http

import (
	"errors"
	"github.com/whosonfirst/go-whosonfirst-placetypes"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"html/template"
	_ "log"
	gohttp "net/http"
)

type PointInPolygonHandlerOptions struct {
	Templates        *template.Template
	InitialLatitude  float64
	InitialLongitude float64
	InitialZoom      int
	DataEndpoint     string
}

type PointInPolygonHandlerTemplateVars struct {
	InitialLatitude  float64
	InitialLongitude float64
	InitialZoom      int
	DataEndpoint     string
	Placetypes       []*placetypes.WOFPlacetype
}

func PointInPolygonHandler(spatial_app *app.SpatialApplication, opts *PointInPolygonHandlerOptions) (gohttp.Handler, error) {

	t := opts.Templates.Lookup("pointinpolygon")

	if t == nil {
		return nil, errors.New("Missing pointinpolygon template")
	}

	walker := spatial_app.Walker

	pt_list, err := placetypes.Placetypes()

	if err != nil {
		return nil, err
	}

	fn := func(rsp gohttp.ResponseWriter, req *gohttp.Request) {

		if walker.IsIndexing() {
			gohttp.Error(rsp, "indexing records", gohttp.StatusServiceUnavailable)
			return
		}

		// important if we're trying to use this in a Lambda/API Gateway context

		rsp.Header().Set("Content-Type", "text/html; charset=utf-8")

		vars := PointInPolygonHandlerTemplateVars{
			InitialLatitude:  opts.InitialLatitude,
			InitialLongitude: opts.InitialLongitude,
			InitialZoom:      opts.InitialZoom,
			DataEndpoint:     opts.DataEndpoint,
			Placetypes:       pt_list,
		}

		err := t.Execute(rsp, vars)

		if err != nil {
			gohttp.Error(rsp, err.Error(), gohttp.StatusInternalServerError)
			return
		}

		return
	}

	h := gohttp.HandlerFunc(fn)
	return h, nil
}