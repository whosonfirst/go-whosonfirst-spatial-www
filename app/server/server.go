package server

import (
	"context"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	gohttp "net/http"
	"path/filepath"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/aaronland/go-http-maps/v2"
	"github.com/aaronland/go-http-ping/v2"
	"github.com/aaronland/go-http-server"
	"github.com/rs/cors"
	"github.com/sfomuseum/go-http-auth"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/http"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/http/api"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/http/www"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/templates/html"
	app "github.com/whosonfirst/go-whosonfirst-spatial/application"
)

func Run(ctx context.Context) error {

	fs, err := DefaultFlagSet()

	if err != nil {
		return fmt.Errorf("Failed to derive default flag set, %w", err)
	}

	return RunWithFlagSet(ctx, fs)
}

func RunWithFlagSet(ctx context.Context, fs *flag.FlagSet) error {

	opts, err := RunOptionsFromFlagSet(ctx, fs)

	if err != nil {
		return fmt.Errorf("Failed to derive options from flag set, %w", err)
	}

	return RunWithOptions(ctx, opts)
}

func RunWithOptions(ctx context.Context, opts *RunOptions) error {

	logger := slog.Default()

	spatial_opts := &app.SpatialApplicationOptions{
		SpatialDatabaseURI:     opts.SpatialDatabaseURI,
		PropertiesReaderURI:    opts.PropertiesReaderURI,
		EnableCustomPlacetypes: opts.EnableCustomPlacetypes,
		CustomPlacetypes:       opts.CustomPlacetypes,
	}

	spatial_app, err := app.NewSpatialApplication(ctx, spatial_opts)

	if err != nil {
		return fmt.Errorf("Failed to create new spatial application, %w", err)
	}

	authenticator, err := auth.NewAuthenticator(ctx, opts.AuthenticatorURI)

	if err != nil {
		return fmt.Errorf("Failed to create authenticator, %w", err)
	}

	go func() {

		err := spatial_app.IndexDatabaseWithIterators(ctx, opts.IteratorSources)

		if err != nil {
			slog.Error("Failed to index database with iterator", "error", err)
		}
	}()

	mux := gohttp.NewServeMux()

	ping_handler, err := ping.PingPongHandler()

	if err != nil {
		return fmt.Errorf("failed to create ping handler because %s", err)
	}

	mux.Handle(opts.PathPing, ping_handler)

	var cors_wrapper *cors.Cors

	if opts.EnableCORS {
		cors_wrapper = cors.New(cors.Options{
			AllowedOrigins:   opts.CORSOrigins,
			AllowCredentials: opts.CORSAllowCredentials,
		})
	}

	// data (geojson) handlers
	// SpatialDatabase implements reader.Reader

	data_handler, err := api.NewDataHandler(spatial_app.SpatialDatabase)

	if err != nil {
		return fmt.Errorf("Failed to create data handler, %v", err)
	}

	data_handler = http.CheckIndexingHandler(spatial_app, data_handler)

	data_handler = authenticator.WrapHandler(data_handler)

	if opts.EnableCORS {
		data_handler = cors_wrapper.Handler(data_handler)
	}

	if opts.EnableGzip {
		data_handler = gziphandler.GzipHandler(data_handler)
	}

	if !strings.HasSuffix(opts.PathData, "/") {
		opts.PathData = fmt.Sprintf("%s/", opts.PathData)
	}

	mux.Handle(opts.PathData, data_handler)

	// point-in-polygon handlers

	api_pip_opts := &api.PointInPolygonHandlerOptions{
		EnableGeoJSON: opts.EnableGeoJSON,
		LogTimings:    opts.LogTimings,
	}

	api_pip_handler, err := api.PointInPolygonHandler(spatial_app, api_pip_opts)

	if err != nil {
		return fmt.Errorf("failed to create point-in-polygon handler because %s", err)
	}

	api_pip_handler = authenticator.WrapHandler(api_pip_handler)

	if opts.EnableCORS {
		api_pip_handler = cors_wrapper.Handler(api_pip_handler)
	}

	if opts.EnableGzip {
		api_pip_handler = gziphandler.GzipHandler(api_pip_handler)
	}

	path_api_pip := filepath.Join(opts.PathAPI, "point-in-polygon")

	mux.Handle(path_api_pip, api_pip_handler)

	// www handlers

	if opts.EnableWWW {

		maps_opts := &maps.AssignMapConfigHandlerOptions{
			MapProvider:       opts.MapProvider,
			MapTileURI:        opts.MapTileURI,
			InitialView:       opts.InitialView,
			LeafletStyle:      opts.LeafletStyle,
			LeafletPointStyle: opts.LeafletPointStyle,
			ProtomapsTheme:    opts.ProtomapsTheme,
		}

		err := maps.AssignMapConfigHandler(maps_opts, mux, "/map.json")

		if err != nil {
			return fmt.Errorf("Failed to assign map config handler, %w", err)
		}

		t := template.New("spatial")

		t = t.Funcs(map[string]interface{}{

			"EnsureRoot": func(path string) string {

				path = strings.TrimLeft(path, "/")

				if opts.PathPrefix == "" {
					return "/" + path
				}

				path = filepath.Join(opts.PathPrefix, path)
				return path
			},

			"DataRoot": func() string {

				path := opts.PathData

				if opts.PathPrefix != "" {
					path = filepath.Join(opts.PathPrefix, path)
				}

				return path
			},

			"APIRoot": func() string {

				path := opts.PathAPI

				if opts.PathPrefix != "" {
					path = filepath.Join(opts.PathPrefix, path)
				}

				return path
			},
		})

		t, err = t.ParseFS(html.FS, "*.html")

		if err != nil {
			return fmt.Errorf("Unable to parse templates, %v", err)
		}

		// point-in-polygon page

		http_pip_opts := &www.PointInPolygonHandlerOptions{
			Templates: t,
		}

		http_pip_handler, err := www.PointInPolygonHandler(spatial_app, http_pip_opts)

		if err != nil {
			return fmt.Errorf("failed to create (bundled) www handler because %s", err)
		}

		http_pip_handler = authenticator.WrapHandler(http_pip_handler)

		mux.Handle(opts.PathPIP, http_pip_handler)

		if !strings.HasSuffix(opts.PathPIP, "/") {
			path_pip_slash := fmt.Sprintf("%s/", opts.PathPIP)
			mux.Handle(path_pip_slash, http_pip_handler)
		}

		// index / splash page

		index_opts := &www.IndexHandlerOptions{
			Templates: t,
		}

		index_handler, err := www.IndexHandler(index_opts)

		if err != nil {
			return fmt.Errorf("Failed to create index handler, %v", err)
		}

		index_handler = authenticator.WrapHandler(index_handler)

		path_index := "/"

		mux.Handle(path_index, index_handler)
	}

	s, err := server.NewServer(ctx, opts.ServerURI)

	if err != nil {
		return fmt.Errorf("Failed to create new server for '%s', %v", server_uri, err)
	}

	logger.Info("Listening for requests", "address", s.Address())

	err = s.ListenAndServe(ctx, mux)

	if err != nil {
		return fmt.Errorf("Failed to start server, %v", err)
	}

	return nil
}
