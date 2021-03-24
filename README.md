# go-whosonfirst-spatial-www

Opionated web application for the `go-whosonfirst-spatial` packages.

## IMPORTANT

This is work in progress. Documentation to follow.

## Tools

To build binary versions of these tools run the `cli` Makefile target. For example:

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

```
$> ./bin/server -h
  -custom-placetypes string
    	A JSON-encoded string containing custom placetypes defined using the syntax described in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-cors
    	Enable CORS headers for data-related and API handlers.
  -enable-custom-placetypes
    	Enable wof:placetype values that are not explicitly defined in the whosonfirst/go-whosonfirst-placetypes repository.
  -enable-gzip
    	Enable gzip-encoding for data-related and API handlers.
  -enable-properties
    	Enable support for 'properties' parameters in queries.
  -enable-tangram
    	Use Tangram.js for rendering map tiles
  -enable-www
    	Enable the interactive /debug endpoint to query points and display results.
  -index-properties
    	Index properties reader.
  -initial-latitude float
    	The initial latitude for map views to use. (default 37.616906)
  -initial-longitude float
    	The initial longitude for map views to use. (default -122.386665)
  -initial-zoom int
    	The initial zoom level for map views to use. (default 14)
  -is-wof
    	Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported schemes are: directory://, featurecollection://, file://, filelist://, geojsonl://, repo://. (default "repo://")
  -leaflet-tile-url string
    	A valid Leaflet (slippy map) tile template URL to use for rendering maps (if -enable-tangram is false)
  -nextzen-apikey string
    	A valid Nextzen API key
  -nextzen-style-url string
    	The URL for the style bundle file to use for maps rendered with Tangram.js (default "/tangram/refill-style.zip")
  -nextzen-tile-url string
    	The URL for Nextzen tiles to use for maps rendered with Tangram.js (default "https://{s}.tile.nextzen.org/tilezen/vector/v1/512/all/{z}/{x}/{y}.mvt")
  -path-data string
    	The URL for data (GeoJSON) handler (default "/data")
  -path-ping string
    	The URL for the ping (health check) handler (default "/health/ping")
  -path-pip string
    	The URL for the point in polygon web handler (default "/point-in-polygon")
  -path-prefix string
    	Prepend this prefix to all assets (but not HTTP handlers). This is mostly for API Gateway integrations.
  -path-root-api string
    	The root URL for all API handlers (default "/api")
  -properties-reader-uri string
    	Valid options are: [rtree://] (default "rtree://")
  -server-uri string
    	A valid aaronland/go-http-server URI. (default "http://localhost:8080")
  -spatial-database-uri string
    	Valid options are: [rtree://]
  -verbose
    	Be chatty.
```

For example:

```
$> bin/server \
	-enable-www \
	-spatial-database-uri 'rtree:///?strict=false' \
	-properties-reader-uri 'whosonfirst://?reader=fs:////usr/local/data/sfomuseum-data-architecture/data&cache=gocache://' \
	-enable-tangram \
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
	-properties-reader-uri 'whosonfirst://?reader=fs:////usr/local/data/sfomuseum-data-architecture/data&cache=gocache://' \
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

If you are returning results as a GeoJSON `FeatureCollection` you may also request additional properties be appended by specifying them as a comma-separated list in the `?properties=` parameter. For example:

```
$> curl -H 'Accept: application/geo+json' -XPOST http://localhost:8080/api/point-in-polygon -d '{"latitude": 37.61701894316063, "longitude": -122.3866653442383, "properties": ["sfomuseum:*" ]}'
{
  "type": "FeatureCollection",
  "features": [
    {
      "type": "Feature",
      "geometry": {
        "type": "MultiPolygon",
        "coordinates": [ ... ]
      },
      "properties": {
        "mz:is_ceased": 1,
        "mz:is_current": 0,
        "mz:is_deprecated": 0,
        "mz:is_superseded": 1,
        "mz:is_superseding": 1,
        "mz:latitude": 37.617037,
        "mz:longitude": -122.385975,
        "mz:max_latitude": 37.62120978585632,
        "mz:max_longitude": -122.38125166743595,
        "mz:min_latitude": 37.61220882045874,
        "mz:min_longitude": -122.39033463643914,
        "mz:uri": "https://data.whosonfirst.org/115/939/632/7/1159396327.geojson",
        "sfomuseum:building_id": "SFO",
        "sfomuseum:is_sfo": 1,
        "sfomuseum:placetype": "building",
        "wof:country": "US",
        "wof:id": 1159396327,
        "wof:lastmodified": 1547232162,
        "wof:name": "SFO Terminal Complex",
        "wof:parent_id": 102527513,
        "wof:path": "115/939/632/7/1159396327.geojson",
        "wof:placetype": "building",
        "wof:repo": "sfomuseum-data-architecture",
        "wof:superseded_by": [
          1159554801
        ],
        "wof:supersedes": [
          1159396331
        ]
      }
    }... and so on
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