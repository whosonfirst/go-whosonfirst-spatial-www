package server

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/whosonfirst/go-whosonfirst-spatial/geo"
	"github.com/sfomuseum/go-flags/multi"
	"log"
	"strconv"
	"strings"
)

var path_prefix string

var path_api string
var path_ping string
var path_pip string
var path_data string

var enable_www bool
var enable_geojson bool

var enable_cors bool
var enable_cors_credentials bool

var cors_origins multi.MultiCSVString

var enable_gzip bool

var enable_tangram bool

var nextzen_apikey string
var nextzen_style_url string
var nextzen_tile_url string

var leaflet_tile_url string

var leaflet_initial_latitude float6
var leaflet_initial_longitude float64
var leaflet_initial_zoom int

var max_bounds string

var server_uri string
var authenticator_uri string

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.StringVar(&server_uri, "server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.StringVar(&authenticator_uri, "authenticator-uri", "null://", "A valid sfomuseum/go-http-auth URI.")
	
	fs.BoolVar(&enable_www, "enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.BoolVar(&enable_geojson, "enable-geojson", false, "Enable GeoJSON output for point-in-polygon API calls.")

	fs.BoolVar(&enable_cors, "enable-cors", false, "Enable CORS headers for data-related and API handlers.")
	fs.BoolVar(&enable_cors_credentials, "enable-cors-credentials", false, "...")	

	fs.Var(&cors_origins, "cors-origin", "...")
	
	fs.BoolVar(&enable_gzip, "enable-gzip", false, "Enable gzip-encoding for data-related and API handlers.")

	fs.StringVar(&path_prefix, "path-prefix", "", "Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.")

	fs.StringVar(&path_api, "path-api", "/api", "The root URL for all API handlers")
	fs.StringVar(&path_ping, "path-ping", "/health/ping", "The URL for the ping (health check) handler")
	fs.StringVar(&path_pip, "path-pip", "/point-in-polygon", "The URL for the point in polygon web handler")
	fs.StringVar(PATH_DATA, "/data", "The URL for data (GeoJSON) handler")

	leaflet_desc := fmt.Sprintf("A valid Leaflet (slippy map) tile template URL to use for rendering maps (if -%s is false)", ENABLE_TANGRAM)
	fs.StringVar(LEAFLET_TILE_URL, "", leaflet_desc)

	fs.BoolVar(ENABLE_TANGRAM, false, "Use Tangram.js for rendering map tiles")

	fs.StringVar(NEXTZEN_APIKEY, "", "A valid Nextzen API key")
	fs.StringVar(NEXTZEN_STYLE_URL, "/tangram/refill-style.zip", "The URL for the style bundle file to use for maps rendered with Tangram.js")
	fs.StringVar(NEXTZEN_TILE_URL, tangramjs.NEXTZEN_MVT_ENDPOINT, "The URL for Nextzen tiles to use for maps rendered with Tangram.js")

	fs.FloatVar(INITIAL_LATITUDE, 37.616906, "The initial latitude for map views to use.")
	fs.FloatVar(INITIAL_LONGITUDE, -122.386665, "The initial longitude for map views to use.")
	fs.Int(INITIAL_ZOOM, 14, "The initial zoom level for map views to use.")
	fs.StringVar(MAX_BOUNDS, "", "An optional comma-separated bounding box ({MINX},{MINY},{MAXX},{MAXY}) to set the boundary for map views.")

	return nil
}

func ValidateWWWFlags(fs *flag.FlagSet) error {

	enable_www, err := lookup.BoolVar(fs, ENABLE_WWW)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", ENABLE_WWW, err)
	}

	if !enable_www {
		return nil
	}

	bool_flags := []string{
		ENABLE_CORS,
		ENABLE_GZIP,
		ENABLE_GEOJSON,
	}

	for _, k := range bool_flags {

		_, err := lookup.BoolVar(fs, k)

		if err != nil {
			return fmt.Errorf("Failed to lookup %s flag, %w", k, err)
		}
	}

	init_lat, err := lookup.Float64Var(fs, INITIAL_LATITUDE)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", INITIAL_LATITUDE, err)
	}

	if !geo.IsValidLatitude(init_lat) {
		return fmt.Errorf("Invalid latitude for %s flag", INITIAL_LATITUDE)
	}

	init_lon, err := lookup.Float64Var(fs, INITIAL_LONGITUDE)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", INITIAL_LONGITUDE, err)
	}

	if !geo.IsValidLongitude(init_lon) {
		return fmt.Errorf("Invalid longitude for %s flag", INITIAL_LONGITUDE)
	}

	init_zoom, err := lookup.IntVar(fs, INITIAL_ZOOM)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", INITIAL_ZOOM, err)
	}

	if init_zoom < 1 {
		return fmt.Errorf("Invalid zoom for %s flag", INITIAL_ZOOM)
	}

	max_bounds, err := lookup.StringVar(fs, MAX_BOUNDS)

	if max_bounds != "" {

		bounds := strings.Split(max_bounds, ",")

		if len(bounds) != 4 {
			return fmt.Errorf("Invalid max bounds for %s flag; expected 4 parts but got %d", MAX_BOUNDS, len(bounds))
		}

		minx, err := strconv.ParseFloat(bounds[0], 64)

		if err != nil {
			return fmt.Errorf("Invalid minx for %s flag, %w", MAX_BOUNDS, err)
		}

		if !geo.IsValidLongitude(minx) {
			return fmt.Errorf("Invalid minx (longitude) for %s flag", MAX_BOUNDS)
		}

		miny, err := strconv.ParseFloat(bounds[1], 64)

		if err != nil {
			return fmt.Errorf("Invalid miny for %s flag, %w", MAX_BOUNDS, err)
		}

		if !geo.IsValidLatitude(miny) {
			return fmt.Errorf("Invalid miny (latitude) for %s flag", MAX_BOUNDS)
		}

		maxx, err := strconv.ParseFloat(bounds[2], 64)

		if err != nil {
			return fmt.Errorf("Invalid maxx for %s flag, %w", MAX_BOUNDS, err)
		}

		if !geo.IsValidLongitude(maxx) {
			return fmt.Errorf("Invalid maxx (longitude) for %s flag", MAX_BOUNDS)
		}

		maxy, err := strconv.ParseFloat(bounds[3], 64)

		if err != nil {
			return fmt.Errorf("Invalid maxy for %s flag, %w", MAX_BOUNDS, err)
		}

		if !geo.IsValidLatitude(maxy) {
			return fmt.Errorf("Invalid maxy (latitude) for %s flag", MAX_BOUNDS)
		}
	}

	path_flags := []string{
		PATH_PREFIX,
		PATH_API,
		PATH_DATA,
		PATH_PING,
		PATH_PIP,
	}

	for _, fl := range path_flags {

		_, err := lookup.StringVar(fs, fl)

		if err != nil {
			return fmt.Errorf("Failed to lookup %s flag, %w", fl, err)
		}
	}

	enable_tangram, err := lookup.BoolVar(fs, ENABLE_TANGRAM)

	if err != nil {
		return fmt.Errorf("Failed to lookup %s flag, %w", ENABLE_TANGRAM, err)
	}

	if enable_tangram {

		nz_keys := []string{
			NEXTZEN_APIKEY,
			NEXTZEN_STYLE_URL,
			NEXTZEN_TILE_URL,
		}

		for _, k := range nz_keys {

			v, err := lookup.StringVar(fs, k)

			if err != nil {
				return fmt.Errorf("Failed to lookup %s flag, %w", k, err)
			}

			if v == "" {
				log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", k)
			}
		}

	} else {

		v, err := lookup.StringVar(fs, LEAFLET_TILE_URL)

		if err != nil {
			return fmt.Errorf("Failed to lookup %s flag, %w", LEAFLET_TILE_URL, err)
		}

		if v == "" {
			log.Printf("-%s flag is empty, this will probably result in unexpected behaviour\n", LEAFLET_TILE_URL)
		}
	}

	return nil
}
