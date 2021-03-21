package api

import (
	"encoding/json"
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-reader"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/output"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/parameters"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spr/v2"
	"github.com/whosonfirst/go-whosonfirst-spr-geojson"
	_ "log"
	"net/http"
)

type PointInPolygonHandlerOptions struct {
	EnableGeoJSON    bool
	EnableProperties bool
	GeoJSONReader    reader.Reader
	SPRPathResolver  geojson.SPRPathResolver
}

func PointInPolygonHandler(spatial_app *app.SpatialApplication, opts *PointInPolygonHandlerOptions) (http.Handler, error) {

	spatial_db := spatial_app.SpatialDatabase
	properties_r := spatial_app.PropertiesReader
	walker := spatial_app.Walker

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()
		query := req.URL.Query()

		coord, err := parameters.Coordinate(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		str_format, err := sanitize.GetString(req, "format")

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		if str_format == "geojson" && !opts.EnableGeoJSON {
			http.Error(rsp, "GeoJSON formatting is disabled.", http.StatusBadRequest)
			return
		}

		if str_format == "properties" && !opts.EnableProperties {
			http.Error(rsp, "Properties formatting is disabled.", http.StatusBadRequest)
			return
		}

		properties_paths, err := parameters.Properties(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		filters, err := filter.NewSPRFilterFromQuery(query)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		results, err := spatial_db.PointInPolygon(ctx, coord, filters)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		if results == nil {
			http.Error(rsp, "Unable to yield results", http.StatusInternalServerError)
			return
		}

		rsp.Header().Set("Content-Type", "application/json")

		var final interface{}
		final = results

		enc := json.NewEncoder(rsp)

		switch str_format {
		case "geojson":

			as_opts := &geojson.AsFeatureCollectionOptions{
				Reader:          opts.GeoJSONReader,
				Writer:          rsp,
				SPRPathResolver: opts.SPRPathResolver,
			}

			err := geojson.AsFeatureCollection(ctx, results, as_opts)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			return

		case "properties":

			if len(properties_paths) > 0 {

				props, err := properties_r.PropertiesResponseResultsWithStandardPlacesResults(ctx, final.(spr.StandardPlacesResults), properties_paths)

				if err != nil {
					http.Error(rsp, err.Error(), http.StatusInternalServerError)
					return
				}

				final = props
			}

		default:
			// spr (above)
		}

		err = enc.Encode(final)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	h := http.HandlerFunc(fn)
	return h, nil
}

func PointInPolygonCandidatesHandler(spatial_app *app.SpatialApplication) (http.Handler, error) {

	walker := spatial_app.Walker
	spatial_db := spatial_app.SpatialDatabase

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		if walker.IsIndexing() {
			http.Error(rsp, "indexing records", http.StatusServiceUnavailable)
			return
		}

		ctx := req.Context()

		coord, err := parameters.Coordinate(req)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusBadRequest)
			return
		}

		candidates, err := spatial_db.PointInPolygonCandidates(ctx, coord)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		output.AsJSON(rsp, candidates)
	}

	h := http.HandlerFunc(fn)
	return h, nil
}
