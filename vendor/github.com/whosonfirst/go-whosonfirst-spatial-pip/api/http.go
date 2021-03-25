package api

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spatial"
	"github.com/whosonfirst/go-whosonfirst-spatial-pip"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	_ "log"
	"net/http"
	"strings"
)

const GEOJSON string = "application/geo+json"

type PointInPolygonHandlerOptions struct {
	EnableGeoJSON bool
}

func PointInPolygonHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		ctx := req.Context()

		if req.Method != "POST" {
			http.Error(rsp, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		if app.Iterator.IsIndexing() {
			http.Error(rsp, "Indexing records", http.StatusServiceUnavailable)
			return
		}

		var pip_req *pip.PointInPolygonRequest

		dec := json.NewDecoder(req.Body)
		err := dec.Decode(&pip_req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		accept, err := sanitize.HeaderString(req, "Accept")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if accept == GEOJSON && !opts.EnableGeoJSON {
			http.Error(rsp, "GeoJSON output is not supported", http.StatusBadRequest)
			return
		}

		str_props, err := sanitize.HeaderString(req, "X-Properties")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}
				
		pip_rsp, err := pip.QueryPointInPolygon(ctx, app, pip_req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if opts.EnableGeoJSON && accept == GEOJSON {

			opts := &geojson.AsFeatureCollectionOptions{
				Reader: app.SpatialDatabase,
				Writer: rsp,
			}

			err := geojson.AsFeatureCollection(ctx, pip_rsp, opts)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		var props []string

		str_props = strings.Trim(str_props, " ")
		
		if str_props != "" {
			props = strings.Split(str_props, ",")
		}
		
		if len(props) > 0 {

			props_opts := &spatial.PropertiesResponseOptions{
				Reader:       app.PropertiesReader,
				Keys:         props,
				SourcePrefix: "properties",
			}

			props_rsp, err := spatial.PropertiesResponseResultsWithStandardPlacesResults(ctx, props_opts, pip_rsp)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			enc := json.NewEncoder(rsp)
			err = enc.Encode(props_rsp)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		enc := json.NewEncoder(rsp)
		err = enc.Encode(pip_rsp)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	pip_handler := http.HandlerFunc(fn)
	return pip_handler, nil
}
