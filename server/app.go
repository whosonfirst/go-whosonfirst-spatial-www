package server

// This is a first-cut at making the core application in to an extensible
// package - it is likely that it will change (20201207/thisisaaronland)

import (
	"context"
	"flag"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-ping"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial-pip/api"
	www_flags "github.com/whosonfirst/go-whosonfirst-spatial-www/flags"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/http"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/templates/html"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	"github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"html/template"
	"log"
	gohttp "net/http"
	"path/filepath"
)

type HTTPServerApplication struct {
}

func NewHTTPServerApplication(ctx context.Context) (*HTTPServerApplication, error) {

	server_app := &HTTPServerApplication{}
	return server_app, nil
}

func (server_app *HTTPServerApplication) Run(ctx context.Context) error {

	fs, err := flags.CommonFlags()

	if err != nil {
		return fmt.Errorf("Failed to instantiate common flags, %v", err)
	}

	err = flags.AppendIndexingFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to append indexings flags, %v", err)
	}

	err = www_flags.AppendWWWFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to append www flags, %v", err)
	}

	flagset.Parse(fs)

	return server_app.RunWithFlagSet(ctx, fs)
}

func (server_app *HTTPServerApplication) RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	err := flags.ValidateCommonFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate common flags, %v", err)
	}

	err = flags.ValidateIndexingFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate indexing flags, %v", err)
	}

	err = www_flags.ValidateWWWFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate www flags, %v", err)
	}

	enable_www, _ := lookup.BoolVar(fs, "enable-www")
	// enable_candidates, _ := lookup.BoolVar(fs, "enable-candidates")

	enable_geojson, _ := lookup.BoolVar(fs, "enable-geojson")

	nextzen_apikey, _ := lookup.StringVar(fs, "nextzen-apikey")
	nextzen_style_url, _ := lookup.StringVar(fs, "nextzen-style-url")
	nextzen_tile_url, _ := lookup.StringVar(fs, "nextzen-tile-url")

	initial_lat, _ := lookup.Float64Var(fs, "initial-latitude")
	initial_lon, _ := lookup.Float64Var(fs, "initial-longitude")
	initial_zoom, _ := lookup.IntVar(fs, "initial-zoom")

	server_uri, _ := lookup.StringVar(fs, "server-uri")

	data_endpoint, _ := lookup.StringVar(fs, "data-endpoint")

	spatial_app, err := app.NewSpatialApplicationWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new spatial application, because: %v", err))
	}

	logger := spatial_app.Logger

	paths := fs.Args()

	err = spatial_app.IndexPaths(ctx, paths...)

	if err != nil {
		return fmt.Errorf("Failed to index paths, because %s", err)
	}

	// START OF ...

	path_prefix, _ := lookup.StringVar(fs, www_flags.PATH_PREFIX)
	path_ping, _ := lookup.StringVar(fs, www_flags.PATH_PING)

	path_api_pip, _ := lookup.StringVar(fs, www_flags.PATH_API_PIP)
	path_www_pip, _ := lookup.StringVar(fs, www_flags.PATH_WWW_PIP)
	path_www_index, _ := lookup.StringVar(fs, www_flags.PATH_WWW_INDEX)

	if path_prefix != "" {

		path_ping = filepath.Join(path_prefix, path_ping)
		path_api_pip = filepath.Join(path_prefix, path_api_pip)
		path_www_pip = filepath.Join(path_prefix, path_www_pip)
		path_www_index = filepath.Join(path_prefix, path_www_index)

	}

	// END OF ...

	mux := gohttp.NewServeMux()

	ping_handler, err := ping.PingHandler()

	if err != nil {
		return fmt.Errorf("failed to create ping handler because %s", err)
	}

	mux.Handle(path_ping, ping_handler)

	enable_cors := true
	enable_gzip := true

	cors_origins := []string{"*"}

	var cors_wrapper *cors.Cors

	if enable_cors {
		cors_wrapper = cors.New(cors.Options{
			AllowedOrigins: cors_origins,
		})
	}

	api_pip_opts := &api.PointInPolygonHandlerOptions{
		EnableGeoJSON: enable_geojson,
	}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	if enable_cors {
		api_pip_handler = cors_wrapper.Handler(api_pip_handler)
	}

	if enable_gzip {
		api_pip_handler = gziphandler.GzipHandler(api_pip_handler)
	}

	logger.Info("Register %s handler", path_api_pip)
	mux.Handle(path_api_pip, api_pip_handler)

	/*
		if enable_candidates {

			logger.Debug("setting up candidates handler")

			candidates_handler, err := api.PointInPolygonCandidatesHandler(spatial_app)

			if err != nil {
				return fmt.Errorf("failed to create Spatial handler because %s", err)
			}

			if enable_cors {
				candidates_handler = cors_wrapper.Handler(candidates_handler)
			}

			mux.Handle("/api/point-in-polygon/candidates", candidates_handler)
		}
	*/

	if enable_www {

		t, err := template.ParseFS(html.FS, "*.html")

		if err != nil {
			return fmt.Errorf("Unable to parse templates, %v", err)
		}

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		tangramjs_opts := tangramjs.DefaultTangramJSOptions()
		tangramjs_opts.Nextzen.APIKey = nextzen_apikey
		tangramjs_opts.Nextzen.StyleURL = nextzen_style_url
		tangramjs_opts.Nextzen.TileURL = nextzen_tile_url

		err = tangramjs.AppendAssetHandlersWithPrefix(mux, path_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append tangram.js assets, %v", err)
		}

		err = bootstrap.AppendAssetHandlersWithPrefix(mux, path_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append bootstrap assets, %v", err)
		}

		err = http.AppendStaticAssetHandlersWithPrefix(mux, path_prefix)

		if err != nil {
			return fmt.Errorf("Failed to append static assets, %v", err)
		}

		http_pip_opts := &http.PointInPolygonHandlerOptions{
			Templates:        t,
			InitialLatitude:  initial_lat,
			InitialLongitude: initial_lon,
			InitialZoom:      initial_zoom,
			DataEndpoint:     data_endpoint,
		}

		http_pip_handler, err := http.PointInPolygonHandler(spatial_app, http_pip_opts)

		if err != nil {
			return fmt.Errorf("failed to create (bundled) www handler because %s", err)
		}

		http_pip_handler = bootstrap.AppendResourcesHandlerWithPrefix(http_pip_handler, bootstrap_opts, path_prefix)
		http_pip_handler = tangramjs.AppendResourcesHandlerWithPrefix(http_pip_handler, tangramjs_opts, path_prefix)

		logger.Info("Register %s handler", path_www_pip)
		mux.Handle(path_www_pip, http_pip_handler)

		index_opts := &http.IndexHandlerOptions{
			Templates: t,
		}

		index_handler, err := http.IndexHandler(index_opts)

		if err != nil {
			return fmt.Errorf("Failed to create index handler, %v", err)
		}

		index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, path_prefix)

		logger.Info("Register %s handler", path_www_index)
		mux.Handle(path_www_index, index_handler)
	}

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server for '%s', %v", server_uri, err)
	}

	logger.Info("Listening on %s", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
