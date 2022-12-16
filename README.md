# go-whosonfirst-spatial-www

Opinionated web application for the `go-whosonfirst-spatial` packages.

## Documentation

Documentation is incomplete at this time.

## Example

```
package main

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"
)

import (
	"context"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/application/server"
	"log"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := server.Run(ctx, logger)

	if err != nil {
		logger.Fatal(err)
	}
}
```

The default `server` implementation uses an in-memory RTree-based spatial index that needs to be populated when the server is started.

There are also server implementations that use SQLite and Protomaps derived spatial databases:

* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-protomaps

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
$> > ./bin/server -h
  -authenticator-uri string
    	A valid sfomuseum/go-http-auth URI. (default "null://")
  -cors-allow-credentials
    	Allow HTTP credentials to be included in CORS requests.
  -cors-origin value
    	One or more hosts to allow CORS requests from; may be a comma-separated list.
  -custom-placetypes string
    	A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-cors
    	Enable CORS headers for data-related and API handlers.
  -enable-custom-placetypes
    	Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-geojson
    	Enable GeoJSON output for point-in-polygon API calls.
  -enable-gzip
    	Enable gzip-encoding for data-related and API handlers.
  -enable-www
    	Enable the interactive /debug endpoint to query points and display results.
  -is-wof
    	Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/v2 URI. Supported schemes are: directory://, featurecollection://, file://, filelist://, geojsonl://, null://, repo://. (default "repo://")
  -leaflet-enable-draw
    	Enable the Leaflet.Draw plugin.
  -leaflet-enable-fullscreen
    	Enable the Leaflet.Fullscreen plugin.
  -leaflet-enable-hash
    	Enable the Leaflet.Hash plugin. (default true)
  -leaflet-initial-latitude float
    	The initial latitude for map views to use. (default 37.616906)
  -leaflet-initial-longitude float
    	The initial longitude for map views to use. (default -122.386665)
  -leaflet-initial-zoom int
    	The initial zoom level for map views to use. (default 14)
  -leaflet-max-bounds string
    	An optional comma-separated bounding box ({MINX},{MINY},{MAXX},{MAXY}) to set the boundary for map views.
  -leaflet-tile-url string
    	A valid Leaflet 'tileLayer' layer URL. Only necessary if -map-provider is "leaflet".
  -log-timings
    	Emit timing metrics to the application's logger
  -map-provider string
    	The name of the map provider to use. Valid options are: leaflet, protomaps, tangram
  -nextzen-apikey string
    	A valid Nextzen API key. Only necessary if -map-provider is "tangram".
  -nextzen-style-url string
    	A valid URL for loading a Tangram.js style bundle. Only necessary if -map-provider is "tangram". (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	A valid Nextzen tile URL template for loading map tiles. Only necessary if -map-provider is "tangram". (default "https://tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -path-api string
    	The root URL for all API handlers (default "/api")
  -path-data string
    	The URL for data (GeoJSON) handler (default "/data")
  -path-ping string
    	The URL for the ping (health check) handler (default "/health/ping")
  -path-pip string
    	The URL for the point in polygon web handler (default "/point-in-polygon")
  -path-prefix string
    	Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.
  -properties-reader-uri string
    	A valid whosonfirst/go-reader.Reader URI. Available options are: [cachereader:// fs:// null:// repo:// stdin://]
  -protomaps-bucket-uri string
    	The gocloud.dev/blob.Bucket URI where Protomaps tiles are stored. Only necessary if -map-provider is "protomaps" and -protomaps-serve-tiles is true.
  -protomaps-caches-size int
    	The size of the internal Protomaps cache if serving tiles locally. Only necessary if -map-provider is "protomaps" and -protomaps-serve-tiles is true. (default 64)
  -protomaps-database string
    	The name of the Protomaps database to serve tiles from. Only necessary if -map-provider is "protomaps" and -protomaps-serve-tiles is true.
  -protomaps-serve-tiles
    	A boolean flag signaling whether to serve Protomaps tiles locally. Only necessary if -map-provider is "protomaps".
  -protomaps-tile-url string
    	A valid Protomaps .pmtiles URL for loading map tiles. Only necessary if -map-provider is "protomaps". (default "/tiles/")
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -spatial-database-uri string
    	A valid whosonfirst/go-whosonfirst-spatial/data.SpatialDatabase URI. options are: [rtree://]
  -tilezen-enable-tilepack
    	Enable to use of Tilezen MBTiles tilepack for tile-serving. Only necessary if -map-provider is "tangram".
  -tilezen-tilepack-path string
    	The path to the Tilezen MBTiles tilepack to use for serving tiles. Only necessary if -map-provider is "tangram" and -tilezen-enable-tilezen is true.
  -verbose
    	Be chatty.
```

For example:

```
$> bin/server \
	-spatial-database-uri 'rtree:///?strict=false' \
	-enable-www \	
	-map-provider tangram \
	-nextzen-apikey {NEXTZEN_APIKEY} \
	/usr/local/data/sfomuseum-data-architecture
	
11:44:31.902988 [main][index] ERROR 1159157931 failed indexing, (rtreego: improper distance). Strict mode is disabled, so skipping.
11:44:32.073804 [main] STATUS finished indexing in 744.717822ms
```

When you visit `http://localhost:8080` in your web browser you should see something like this:

![](docs/images/server.png)

If you don't need, or want, to expose a user-facing interface simply remove the `-enable-www` and `-nextzen-apikey` flags. For example:

```
$> bin/server \
	-enable-geojson \
	-spatial-database-uri 'rtree:///?strict=false' \
	/usr/local/data/sfomuseum-data-architecture
```

And then to query the point-in-polygon API you would do something like this:

```
$> curl -XPOST http://localhost:8080/api/point-in-polygon -d '{"latitude": 37.61701894316063, "longitude": -122.3866653442383}'

{
  "places": [
    {
      "wof:id": 1360665043,
      "wof:parent_id": -1,
      "wof:name": "Central Parking Garage",
      "wof:placetype": "wing",
      "wof:country": "US",
      "wof:repo": "sfomuseum-data-architecture",
      "wof:path": "136/066/504/3/1360665043.geojson",
      "wof:superseded_by": [],
      "wof:supersedes": [
        1360665035
      ],
      "mz:uri": "https://data.whosonfirst.org/136/066/504/3/1360665043.geojson",
      "mz:latitude": 37.616332,
      "mz:longitude": -122.386047,
      "mz:min_latitude": 37.61498599208708,
      "mz:min_longitude": -122.38779093748578,
      "mz:max_latitude": 37.61767331604971,
      "mz:max_longitude": -122.38429192207244,
      "mz:is_current": 0,
      "mz:is_ceased": 1,
      "mz:is_deprecated": 0,
      "mz:is_superseded": 0,
      "mz:is_superseding": 1,
      "wof:lastmodified": 1547232156
    }
    ... and so on
}    
```

By default, results are returned as a list of ["standard places response"](https://github.com/whosonfirst/go-whosonfirst-spr/) (SPR) elements. You can also return results as a GeoJSON `FeatureCollection` by including a `format=geojson` query parameter. For example:


```
$> curl -H 'Accept: application/geo+json' -XPOST http://localhost:8080/api/point-in-polygon -d '{"latitude": 37.61701894316063, "longitude": -122.3866653442383}'

{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "MultiPolygon",
        "coordinates": [ ...omitted for the sake of brevity ]
      },
      "properties": {
        "mz:is_ceased": 1,
        "mz:is_current": 0,
        "mz:is_deprecated": 0,
        "mz:is_superseded": 0,
        "mz:is_superseding": 1,
        "mz:latitude": 37.616332,
        "mz:longitude": -122.386047,
        "mz:max_latitude": 37.61767331604971,
        "mz:max_longitude": -122.38429192207244,
        "mz:min_latitude": 37.61498599208708,
        "mz:min_longitude": -122.38779093748578,
        "mz:uri": "https://data.whosonfirst.org/136/066/504/3/1360665043.geojson",
        "wof:country": "US",
        "wof:id": 1360665043,
        "wof:lastmodified": 1547232156,
        "wof:name": "Central Parking Garage",
        "wof:parent_id": -1,
        "wof:path": "136/066/504/3/1360665043.geojson",
        "wof:placetype": "wing",
        "wof:repo": "sfomuseum-data-architecture",
        "wof:superseded_by": [],
        "wof:supersedes": [
          1360665035
        ]
      }
    }
    ... and so on
  ]
}  
```

### Indexing "plain old" GeoJSON

There is early support for indexing "plain old" GeoJSON, as in GeoJSON documents that do not following the naming conventions for properties that Who's On First documents use. It is very likely there are still bugs or subtle gotchas.

For example, here's how we could index and serve a GeoJSON FeatureCollection of building footprints:

```
$> bin/server
	-spatial-database-uri 'rtree:///?strict=false' \
	-iterator-uri featurecollection:// \
	/usr/local/data/footprint.geojson
```

And then:

```
$> curl -s -XPOST 'http://localhost:8080/api/point-in-polygon '{"latitude": 37.61686957521345, "longitude": -122.3903158758416}' \

| jq '.["places"][]["spr:id"]'

"1014"
"1031"
"1015"
"1026"
```

Support for returning results in the `properties` or `geojson` format is not available for "plain old" GeoJSON records at this time.

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-rtree
* https://github.com/whosonfirst/go-whosonfirst-spatial-pip
* https://github.com/aaronland/go-http-maps