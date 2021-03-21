package flags

import (
	"flag"
	"github.com/aaronland/go-http-tangramjs"
)

func AppendWWWFlags(fs *flag.FlagSet) error {

	fs.String("server-uri", "http://localhost:8080", "A valid aaronland/go-http-server URI.")

	fs.Bool("enable-geojson", false, "Allow users to request GeoJSON FeatureCollection formatted responses.")
	fs.String("geojson-path-resolver-uri", "", "... Necessary for GeoJSON records that don't follow the Who's On First document conventions.")

	// fs.Bool("enable-candidates", false, "Enable the /candidates endpoint to return candidate bounding boxes (as GeoJSON) for requests.")
	fs.Bool("enable-www", false, "Enable the interactive /debug endpoint to query points and display results.")

	fs.String("geojson-reader-uri", "", "A valid whosonfirst/go-reader.Reader URI. Required if the -enable-geojson or -enable-www flags are set.")

	fs.String("www-path", "/debug", "The URL path for the interactive debug endpoint.")

	fs.String("static-prefix", "", "Prepend this prefix to URLs for static assets.")

	fs.String("nextzen-apikey", "", "A valid Nextzen API key")
	fs.String("nextzen-style-url", "/tangram/refill-style.zip", "...")
	fs.String("nextzen-tile-url", tangramjs.NEXTZEN_MVT_ENDPOINT, "...")

	fs.Float64("initial-latitude", 37.616906, "...")
	fs.Float64("initial-longitude", -122.386665, "...")
	fs.Int("initial-zoom", 13, "...")

	return nil
}
