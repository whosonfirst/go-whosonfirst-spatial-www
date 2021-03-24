package flags

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
)

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.String(SERVER_URI, "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.Bool(ENABLE_WWW, false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.Bool(ENABLE_GEOJSON, false, "Enable GeoJSON output for point-in-polygon API calls.")
	
	fs.Bool(ENABLE_CORS, false, "Enable CORS headers for data-related and API handlers.")
	fs.Bool(ENABLE_GZIP, false, "Enable gzip-encoding for data-related and API handlers.")

	fs.String(PATH_PREFIX, "", "Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.")

	fs.String(PATH_API, "/api", "The root URL for all API handlers")
	fs.String(PATH_PING, "/health/ping", "The URL for the ping (health check) handler")
	fs.String(PATH_PIP, "/point-in-polygon", "The URL for the point in polygon web handler")
	fs.String(PATH_DATA, "/data", "The URL for data (GeoJSON) handler")

	leaflet_desc := fmt.Sprintf("A valid Leaflet (slippy map) tile template URL to use for rendering maps (if -%s is false)", ENABLE_TANGRAM)
	fs.String(LEAFLET_TILE_URL, "", leaflet_desc)

	fs.Bool(ENABLE_TANGRAM, false, "Use Tangram.js for rendering map tiles")

	fs.String(NEXTZEN_APIKEY, "", "A valid Nextzen API key")
	fs.String(NEXTZEN_STYLE_URL, "/tangram/refill-style.zip", "The URL for the style bundle file to use for maps rendered with Tangram.js")
	fs.String(NEXTZEN_TILE_URL, tangramjs.NEXTZEN_MVT_ENDPOINT, "The URL for Nextzen tiles to use for maps rendered with Tangram.js")

	fs.Float64(INITIAL_LATITUDE, 37.616906, "The initial latitude for map views to use.")
	fs.Float64(INITIAL_LONGITUDE, -122.386665, "The initial longitude for map views to use.")
	fs.Int(INITIAL_ZOOM, 14, "The initial zoom level for map views to use.")

	return nil
}
