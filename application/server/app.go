package server

// This is a first-cut at making the core application in to an extensible
// package - it is likely that it will change (20201207/thisisaaronland)

import (
	"context"
	"flag"
	"fmt"
	"github.com/NYTimes/gziphandler"
	"github.com/aaronland/go-http-bootstrap"
	"github.com/aaronland/go-http-leaflet"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/sfomuseum/go-http-auth"	
	"github.com/whosonfirst/go-whosonfirst-spatial-pip/api"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/http"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/templates/html"
	"github.com/whosonfirst/go-whosonfirst-spatial/app"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
	"html/template"
	"log"
	gohttp "net/http"
	"path/filepath"
	"strings"
)

func DefaultFlagSet() (*flag.FlagSet, error) {

	fs, err := spatial_flags.CommonFlags()

	if err != nil {
		return nil, fmt.Errorf("Failed to derive common spatial flags, %w", err)
	}

	err = spatial_flags.AppendIndexingFlags(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to append spatial indexing flags, %w", err)
	}

	err = AppendWWWFlags(fs)

	if err != nil {
		return nil, fmt.Errorf("Failed to append www flags, %w", err)
	}

	return fs, nil
}

func Run(ctx context.Context, logger *log.Logger) error {

	fs, err := DefaultFlagSet()

	if err != nil {
		return fmt.Errorf("Failed to derive default flag set, %w", err)
	}

	return RunWithFlagSet(ctx, fs, logger)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet, logger *log.Logger) error {

	flagset.Parse(fs)

	err := flagset.SetFlagsFromEnvVarsWithFeedback(fs, "WHOSONFIRST", true)

	if err != nil {
		return fmt.Errorf("Failed to set flags from environment variables, %v", err)
	}

	err = spatial_flags.ValidateCommonFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate common flags, %v", err)
	}

	err = spatial_flags.ValidateIndexingFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate indexing flags, %v", err)
	}

	err = ValidateWWWFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to validate www flags, %v", err)
	}

	enable_www, _ := lookup.BoolVar(fs, ENABLE_WWW)
	enable_cors, _ := lookup.BoolVar(fs, ENABLE_CORS)
	enable_gzip, _ := lookup.BoolVar(fs, ENABLE_GZIP)
	enable_geojson, _ := lookup.BoolVar(fs, ENABLE_GEOJSON)

	enable_tangram, _ := lookup.BoolVar(fs, ENABLE_TANGRAM)

	nextzen_apikey, _ := lookup.StringVar(fs, NEXTZEN_APIKEY)
	nextzen_style_url, _ := lookup.StringVar(fs, NEXTZEN_STYLE_URL)
	nextzen_tile_url, _ := lookup.StringVar(fs, NEXTZEN_TILE_URL)

	leaflet_tile_url, _ := lookup.StringVar(fs, LEAFLET_TILE_URL)

	initial_lat, _ := lookup.Float64Var(fs, INITIAL_LATITUDE)
	initial_lon, _ := lookup.Float64Var(fs, INITIAL_LONGITUDE)
	initial_zoom, _ := lookup.IntVar(fs, INITIAL_ZOOM)
	max_bounds, _ := lookup.StringVar(fs, MAX_BOUNDS)

	server_uri, _ := lookup.StringVar(fs, SERVER_URI)
	authenticator_uri, _ := lookup.StringVar(fs, AUTHENTICATOR_URI)
	
	spatial_app, err := app.NewSpatialApplicationWithFlagSet(ctx, fs)

	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to create new spatial application, because: %v", err))
	}

	spatial_app.Logger = logger

	authenticator, err := auth.NewAuthenticator(ctx, authenticator_uri)

	if err != nil {
		return fmt.Errorf("Failed to create authenticator, %w", err)
	}
	
	paths := fs.Args()

	go func() {

		err = spatial_app.IndexPaths(ctx, paths...)

		if err != nil {
			log.Printf("Failed to index paths, because %s", err)
		}
	}()

	path_api, _ := lookup.StringVar(fs, PATH_API)
	path_pip, _ := lookup.StringVar(fs, PATH_PIP)
	path_ping, _ := lookup.StringVar(fs, PATH_PING)
	path_data, _ := lookup.StringVar(fs, PATH_DATA)

	path_prefix, _ := lookup.StringVar(fs, PATH_PREFIX)

	mux := gohttp.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("failed to create ping handler because %s", err)
	}

	mux.Handle(path_ping, ping_handler)

	cors_origins := []string{"*"}

	var cors_wrapper *cors.Cors

	if enable_cors {
		cors_wrapper = cors.New(cors.Options{
			AllowedOrigins: cors_origins,
		})
	}

	// data (geojson) handlers
	// SpatialDatabase implements reader.Reader

	data_handler, err := http.NewDataHandler(spatial_app.SpatialDatabase)

	if err != nil {
		return fmt.Errorf("Failed to create data handler, %v", err)
	}

	data_handler = authenticator.WrapHandler(data_handler)
	
	if enable_cors {
		data_handler = cors_wrapper.Handler(data_handler)
	}

	if enable_gzip {
		data_handler = gziphandler.GzipHandler(data_handler)
	}

	if !strings.HasSuffix(path_data, "/") {
		path_data = fmt.Sprintf("%s/", path_data)
	}

	logger.Printf("Register %s handler\n", path_data)
	mux.Handle(path_data, data_handler)

	// point-in-polygon handlers

	api_pip_opts := &api.PointInPolygonHandlerOptions{
		EnableGeoJSON: enable_geojson,
	}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	api_pip_handler = authenticator.WrapHandler(api_pip_handler)
	
	if enable_cors {
		api_pip_handler = cors_wrapper.Handler(api_pip_handler)
	}

	if enable_gzip {
		api_pip_handler = gziphandler.GzipHandler(api_pip_handler)
	}

	path_api_pip := filepath.Join(path_api, "point-in-polygon")

	logger.Printf("Register %s handler\n", path_api_pip)
	mux.Handle(path_api_pip, api_pip_handler)

	// www handlers

	if enable_www {

		t := template.New("spatial")

		t = t.Funcs(map[string]interface{}{

			"EnsureRoot": func(path string) string {

				path = strings.TrimLeft(path, "/")

				if path_prefix == "" {
					return "/" + path
				}

				path = filepath.Join(path_prefix, path)
				return path
			},

			"DataRoot": func() string {

				path := path_data

				if path_prefix != "" {
					path = filepath.Join(path_prefix, path)
				}

				return path
			},

			"APIRoot": func() string {

				path := path_api

				if path_prefix != "" {
					path = filepath.Join(path_prefix, path)
				}

				return path
			},
		})

		t, err := t.ParseFS(html.FS, "*.html")

		if err != nil {
			return fmt.Errorf("Unable to parse templates, %v", err)
		}

		bootstrap_opts := bootstrap.DefaultBootstrapOptions()

		tangramjs_opts := tangramjs.DefaultTangramJSOptions()
		tangramjs_opts.NextzenOptions.APIKey = nextzen_apikey
		tangramjs_opts.NextzenOptions.StyleURL = nextzen_style_url
		tangramjs_opts.NextzenOptions.TileURL = nextzen_tile_url

		tangramjs_opts.LeafletOptions.EnableHash()

		leaflet_opts := leaflet.DefaultLeafletOptions()
		leaflet_opts.EnableHash()

		if enable_tangram {

			err = tangramjs.AppendAssetHandlers(mux)

			if err != nil {
				return fmt.Errorf("Failed to append tangram.js assets, %v", err)
			}

		} else {

			err = leaflet.AppendAssetHandlers(mux)

			if err != nil {
				return fmt.Errorf("Failed to append leaflet.js assets, %v", err)
			}
		}

		err = bootstrap.AppendAssetHandlers(mux)

		if err != nil {
			return fmt.Errorf("Failed to append bootstrap assets, %v", err)
		}

		err = http.AppendStaticAssetHandlers(mux)

		if err != nil {
			return fmt.Errorf("Failed to append static assets, %v", err)
		}

		http_pip_opts := &http.PointInPolygonHandlerOptions{
			Templates:        t,
			InitialLatitude:  initial_lat,
			InitialLongitude: initial_lon,
			InitialZoom:      initial_zoom,
			MaxBounds:        max_bounds,
			LeafletTileURL:   leaflet_tile_url,
		}

		http_pip_handler, err := http.PointInPolygonHandler(spatial_app, http_pip_opts)

		if err != nil {
			return fmt.Errorf("failed to create (bundled) www handler because %s", err)
		}

		http_pip_handler = bootstrap.AppendResourcesHandlerWithPrefix(http_pip_handler, bootstrap_opts, path_prefix)

		if enable_tangram {
			http_pip_handler = tangramjs.AppendResourcesHandlerWithPrefix(http_pip_handler, tangramjs_opts, path_prefix)
		} else {
			http_pip_handler = leaflet.AppendResourcesHandlerWithPrefix(http_pip_handler, leaflet_opts, path_prefix)
		}

		http_pip_handler = authenticator.WrapHandler(http_pip_handler)
		
		logger.Printf("Register %s handler\n", path_pip)
		mux.Handle(path_pip, http_pip_handler)

		if !strings.HasSuffix(path_pip, "/") {
			path_pip_slash := fmt.Sprintf("%s/", path_pip)
			mux.Handle(path_pip_slash, http_pip_handler)
		}

		index_opts := &http.IndexHandlerOptions{
			Templates: t,
		}

		index_handler, err := http.IndexHandler(index_opts)

		if err != nil {
			return fmt.Errorf("Failed to create index handler, %v", err)
		}

		index_handler = bootstrap.AppendResourcesHandlerWithPrefix(index_handler, bootstrap_opts, path_prefix)

		index_handler = authenticator.WrapHandler(index_handler)
		
		path_index := "/"

		logger.Printf("Register %s handler\n", path_index)
		mux.Handle(path_index, index_handler)
	}

	s, err := server.NewServer(ctx, server_uri)

	if err != nil {
		return fmt.Errorf("Failed to create new server for '%s', %v", server_uri, err)
	}

	logger.Printf("Listening on %s\n", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
