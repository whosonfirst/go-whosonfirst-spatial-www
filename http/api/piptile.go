package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aaronland/go-http-sanitize"
	orb_maptile "github.com/paulmach/orb/maptile"
	spatial_app "github.com/whosonfirst/go-whosonfirst-spatial/application"
	"github.com/whosonfirst/go-whosonfirst-spatial/maptile"
	"github.com/whosonfirst/go-whosonfirst-spatial/query"
)

type PointInPolygonTileHandlerOptions struct{}

func PointInPolygonTileHandler(app *spatial_app.SpatialApplication, opts *PointInPolygonTileHandlerOptions) (http.Handler, error) {

	fn := func(rsp http.ResponseWriter, req *http.Request) {

		logger := slog.Default()

		ctx := req.Context()

		z, err := sanitize.GetInt(req, "z")

		if err != nil {
			logger.Error("Failed to derive z", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		x, err := sanitize.GetInt(req, "x")

		if err != nil {
			logger.Error("Failed to derive x", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		y, err := sanitize.GetInt(req, "y")
		if err != nil {
			logger.Error("Failed to derive y", "error", err)
			http.Error(rsp, "Bad request", http.StatusBadRequest)
			return
		}

		zm := orb_maptile.Zoom(uint32(z))
		map_t := orb_maptile.New(uint32(x), uint32(y), zm)

		spatial_q := &query.SpatialQuery{}

		fc, err := maptile.PointInPolygonCandidateFeaturessFromTile(ctx, app.SpatialDatabase, spatial_q, map_t)

		if err != nil {
			return
		}

		rsp.Header().Set("Content-type", GEOJSON)

		enc := json.NewEncoder(rsp)
		err = enc.Encode(fc)

		if err != nil {
			http.Error(rsp, err.Error(), http.StatusInternalServerError)
			return
		}

		return
	}

	piptile_handler := http.HandlerFunc(fn)
	return piptile_handler, nil
}
