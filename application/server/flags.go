package server

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/sfomuseum/go-flags/multi"
	spatial_flags "github.com/whosonfirst/go-whosonfirst-spatial/flags"
)

var path_prefix string

var path_api string
var path_ping string
var path_pip string
var path_data string

var enable_www bool
var enable_geojson bool

var enable_cors bool
var cors_allow_credentials bool

var cors_origins multi.MultiCSVString

var enable_gzip bool

var enable_tangram bool

var nextzen_apikey string
var nextzen_style_url string
var nextzen_tile_url string

var leaflet_tile_url string

var leaflet_initial_latitude float64
var leaflet_initial_longitude float64
var leaflet_initial_zoom int

var leaflet_max_bounds string

var server_uri string
var authenticator_uri string

var log_timings bool

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

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "A valid sfomuseum/go-http-auth URI.")

	fs.BoolVar(&enable_www, "enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.BoolVar(&enable_geojson, "enable-geojson", false, "Enable GeoJSON output for point-in-polygon API calls.")

	fs.BoolVar(&enable_cors, "enable-cors", false, "Enable CORS headers for data-related and API handlers.")
	fs.BoolVar(&cors_allow_credentials, "cors-allow-credentials", false, "...")

	fs.Var(&cors_origins, "cors-origin", "...")

	fs.BoolVar(&enable_gzip, "enable-gzip", false, "Enable gzip-encoding for data-related and API handlers.")

	fs.StringVar(&path_prefix, "path-prefix", "", "Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.")

	fs.StringVar(&path_api, "path-api", "/api", "The root URL for all API handlers")
	fs.StringVar(&path_ping, "path-ping", "/health/ping", "The URL for the ping (health check) handler")
	fs.StringVar(&path_pip, "path-pip", "/point-in-polygon", "The URL for the point in polygon web handler")
	fs.StringVar(&path_data, "path-data", "/data", "The URL for data (GeoJSON) handler")

	leaflet_desc := fmt.Sprintf("A valid Leaflet (slippy map) tile template URL to use for rendering maps (if -%s is false)", "enable-tangram")
	fs.StringVar(&leaflet_tile_url, "leaflet-tile-url", "", leaflet_desc)

	fs.BoolVar(&enable_tangram, "enable-tangram", false, "Use Tangram.js for rendering map tiles")

	fs.StringVar(&nextzen_apikey, "nextzen-apikey", "", "A valid Nextzen API key")
	fs.StringVar(&nextzen_style_url, "nextzen-style-url", "/tangram/refill-style.zip", "The URL for the style bundle file to use for maps rendered with Tangram.js")
	fs.StringVar(&nextzen_tile_url, "nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "The URL for Nextzen tiles to use for maps rendered with Tangram.js")

	fs.Float64Var(&leaflet_initial_latitude, "leaflet-initial-latitude", 37.616906, "The initial latitude for map views to use.")
	fs.Float64Var(&leaflet_initial_longitude, "leaflet-initial-longitude", -122.386665, "The initial longitude for map views to use.")
	fs.IntVar(&leaflet_initial_zoom, "leaflet-initial-zoom", 14, "The initial zoom level for map views to use.")
	fs.StringVar(&leaflet_max_bounds, "leaflet-max-bounds", "", "An optional comma-separated bounding box ({MINX},{MINY},{MAXX},{MAXY}) to set the boundary for map views.")

	fs.BoolVar(&log_timings, "log-timings", false, "Emit timing metrics to the application's logger")
	return nil
}
