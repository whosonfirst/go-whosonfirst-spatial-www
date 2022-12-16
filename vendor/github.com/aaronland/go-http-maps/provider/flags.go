package provider

import (
	"flag"
	"fmt"
	"github.com/aaronland/go-http-tangramjs"
	"net/url"
	"strconv"
	"strings"
)

const MapProviderFlag string = "map-provider"

var map_provider string

const LeafletEnableHashFlag string = "leaflet-enable-hash"

var leaflet_enable_hash bool

const LeafletEnableFullscreenFlag string = "leaflet-enable-fullscreen"

var leaflet_enable_fullscreen bool

const LeafletEnableDrawFlag string = "leaflet-enable-draw"

var leaflet_enable_draw bool

const LeafletTileURLFlag string = "leaflet-tile-url"

var leaflet_tile_url string

const NextzenAPIKeyFlag string = "nextzen-apikey"

var nextzen_apikey string

const NextzenStyleURLFlag string = "nextzen-style-url"

var nextzen_style_url string

const NextzenTileURLFlag string = "nextzen-tile-url"

var nextzen_tile_url string

const ProtomapsTileURLFlag string = "protomaps-tile-url"

var protomaps_tile_url string

const TilezenEnableTilepack string = "tilezen-enable-tilepack"

var tilezen_enable_tilepack bool

const TilezenTilepackPath string = "tilezen-tilepack-path"

var tilezen_tilepack_path string

const ProtomapsServeTilesFlag string = "protomaps-serve-tiles"

var protomaps_serve_tiles bool

const ProtomapsCacheSizeFlag string = "protomaps-caches-size"

var protomaps_cache_size int

const ProtomapsBucketURIFlag string = "protomaps-bucket-uri"

var protomaps_bucket_uri string

const ProtomapsDatabaseFlag string = "protomaps-database"

var protomaps_database string

func AppendProviderFlags(fs *flag.FlagSet) error {

	schemes := Schemes()
	labels := make([]string, len(schemes))

	for idx, s := range schemes {
		labels[idx] = strings.Replace(s, "://", "", 1)
	}

	str_schemes := strings.Join(labels, ", ")
	map_provider_desc := fmt.Sprintf("Valid options are: %s", str_schemes)

	fs.StringVar(&map_provider, MapProviderFlag, "", map_provider_desc)

	err := AppendLeafletFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to append Leaflet flags, %w", err)
	}

	err = AppendTangramProviderFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to append TangramJS flags, %w", err)
	}

	err = AppendProtomapsProviderFlags(fs)

	if err != nil {
		return fmt.Errorf("Failed to append Protomaps flags, %w", err)
	}

	return nil
}

func AppendLeafletFlags(fs *flag.FlagSet) error {

	fs.BoolVar(&leaflet_enable_hash, LeafletEnableHashFlag, true, "Enable the Leaflet.Hash plugin.")
	fs.BoolVar(&leaflet_enable_fullscreen, LeafletEnableFullscreenFlag, false, "Enable the Leaflet.Fullscreen plugin.")
	fs.BoolVar(&leaflet_enable_draw, LeafletEnableDrawFlag, false, "Enable the Leaflet.Draw plugin.")

	fs.StringVar(&leaflet_tile_url, LeafletTileURLFlag, "", "A valid Leaflet tile URL. Only necessary if -map-provider is \"leaflet\".")
	return nil
}

func AppendTangramProviderFlags(fs *flag.FlagSet) error {

	fs.StringVar(&nextzen_apikey, NextzenAPIKeyFlag, "", "A valid Nextzen API key. Only necessary if -map-provider is \"tangram\".")
	fs.StringVar(&nextzen_style_url, NextzenStyleURLFlag, "/tangram/refill-style.zip", "A valid URL for loading a Tangram.js style bundle. Only necessary if -map-provider is \"tangram\".")
	fs.StringVar(&nextzen_tile_url, NextzenTileURLFlag, tangramjs.NEXTZEN_MVT_ENDPOINT, "A valid Nextzen tile URL template for loading map tiles. Only necessary if -map-provider is \"tangram\".")

	fs.BoolVar(&tilezen_enable_tilepack, TilezenEnableTilepack, false, "Enable to use of Tilezen MBTiles tilepack for tile-serving. Only necessary if -map-provider is \"tangram\".")
	fs.StringVar(&tilezen_tilepack_path, TilezenTilepackPath, "", "The path to the Tilezen MBTiles tilepack to use for serving tiles. Only necessary if -map-provider is \"tangram\" and -tilezen-enable-tilezen is true.")

	return nil
}

func AppendProtomapsProviderFlags(fs *flag.FlagSet) error {

	fs.StringVar(&protomaps_tile_url, ProtomapsTileURLFlag, "/tiles/", "A valid Protomaps .pmtiles URL for loading map tiles. Only necessary if -map-provider is \"protomaps\".")

	fs.BoolVar(&protomaps_serve_tiles, ProtomapsServeTilesFlag, false, "A boolean flag signaling whether to serve Protomaps tiles locally. Only necessary if -map-provider is \"protomaps\".")
	fs.IntVar(&protomaps_cache_size, ProtomapsCacheSizeFlag, 64, "The size of the internal Protomaps cache if serving tiles locally. Only necessary if -map-provider is \"protomaps\" and -protomaps-serve-tiles is true.")
	fs.StringVar(&protomaps_bucket_uri, ProtomapsBucketURIFlag, "", "The gocloud.dev/blob.Bucket URI where Protomaps tiles are stored. Only necessary if -map-provider is \"protomaps\" and -protomaps-serve-tiles is true.")
	fs.StringVar(&protomaps_database, ProtomapsDatabaseFlag, "", "The name of the Protomaps database to serve tiles from. Only necessary if -map-provider is \"protomaps\" and -protomaps-serve-tiles is true.")

	return nil
}

func ProviderURIFromFlagSet(fs *flag.FlagSet) (string, error) {

	u := url.URL{}
	u.Scheme = map_provider

	q := url.Values{}

	if leaflet_enable_hash {
		q.Set("leaflet-enable-hash", strconv.FormatBool(leaflet_enable_hash))
	}

	if leaflet_enable_fullscreen {
		q.Set("leaflet-enable-fullscreen", strconv.FormatBool(leaflet_enable_fullscreen))
	}

	if leaflet_enable_draw {
		q.Set("leaflet-enable-draw", strconv.FormatBool(leaflet_enable_draw))
	}

	switch map_provider {
	case "leaflet":

		q.Set(LeafletTileURLFlag, leaflet_tile_url)

	case "protomaps":

		q.Set(ProtomapsTileURLFlag, protomaps_tile_url)

		if protomaps_serve_tiles {

			q.Set(ProtomapsServeTilesFlag, strconv.FormatBool(protomaps_serve_tiles))
			q.Set(ProtomapsCacheSizeFlag, strconv.Itoa(protomaps_cache_size))
			q.Set(ProtomapsBucketURIFlag, protomaps_bucket_uri)
			q.Set(ProtomapsDatabaseFlag, protomaps_database)
		}

	case "tangram":

		q.Set("nextzen-apikey", nextzen_apikey)

		if nextzen_style_url != "" {
			q.Set("nextzen-style-url", nextzen_style_url)
		}

		if nextzen_tile_url != "" {
			q.Set("nextzen-tile-url", nextzen_tile_url)
		}

		if tilezen_enable_tilepack {

			q.Set("tilezen-enable-tilepack", strconv.FormatBool(tilezen_enable_tilepack))
			q.Set("tilezen-tilepack-path", tilezen_tilepack_path)
			q.Set("tilezen-tilepack-url", "/tilezen/")

			q.Del("nextzen-tile-url")
			q.Set("nextzen-tile-url", "/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
		}

	default:
		// pass
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}
