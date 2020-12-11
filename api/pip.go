package api

import (
	"github.com/aaronland/go-http-sanitize"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/output"
	"github.com/whosonfirst/go-whosonfirst-spatial-http/api/parameters"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/filter"
	"github.com/whosonfirst/go-whosonfirst-spr"
	_ "log"
	"net/http"
)

type PointInPolygonHandlerOptions struct {
	EnableGeoJSON    bool
	EnableProperties bool
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
			http.Error(rsp, "Invalid format", http.StatusBadRequest)
			return
		}

		if str_format == "properties" && !opts.EnableProperties {
			http.Error(rsp, "Invalid format", http.StatusBadRequest)
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
		
		var final interface{}
		final = results

		switch str_format {
		case "geojson":

			collection, err := spatial_db.StandardPlacesResultsToFeatureCollection(ctx, results)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			if collection == nil {
				http.Error(rsp, "Unable to yield GeoJSON results", http.StatusInternalServerError)
				return				
			}
			
			err = properties_r.AppendPropertiesWithFeatureCollection(ctx, collection, properties_paths)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			final = collection

		case "properties":

			props, err := properties_r.PropertiesResponseResultsWithStandardPlacesResults(ctx, final.(spr.StandardPlacesResults), properties_paths)

			if err != nil {
				http.Error(rsp, err.Error(), http.StatusInternalServerError)
				return
			}

			final = props

		default:
			// spr (above)
		}

		output.AsJSON(rsp, final)
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
