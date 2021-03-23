package flags

import (
	"flag"
	"github.com/aaronland/go-http-tangramjs"
)

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.Bool("enable-geojson", false, "Allow users to request GeoJSON FeatureCollection formatted responses. This is automatically enabled if the -enable-www flag is set.")

	fs.Bool("enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.String(PATH_PREFIX, "", "Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.")

	fs.String(PATH_API, "/api", "The root URL for all API handlers")
	fs.String(PATH_PING, "/health/ping", "The URL for the ping (health check) handler")
	fs.String(PATH_PIP, "/point-in-polygon", "The URL for the point in polygon web handler")

	fs.String("nextzen-apikey", "", "A valid Nextzen API key")
	fs.String("nextzen-style-url", "/tangram/refill-style.zip", "...")
	fs.String("nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "...")

	fs.Float64("initial-latitude", 37.616906, "...")
	fs.Float64("initial-longitude", -122.386665, "...")
	fs.Int("initial-zoom", 14, "...")

	return nil
}
