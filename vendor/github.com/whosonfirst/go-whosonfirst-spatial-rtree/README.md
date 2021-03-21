# go-whosonfirst-spatial-rtree

In-memory implementation of the go-whosonfirst-spatial interfaces.

## Important

This is work in progress. Document remains incomplete.

## Interfaces

This package implements the following [go-whosonfirst-spatial](#) interfaces.

### spatial.SpatialDatabase

```
import (
	"github.com/whosonfirst/go-whosonfirst-spatial/database"
	_ "github.com/whosonfirst/go-whosonfirst-spatial-rtree"       
)

db, err := database.NewSpatialDatabase(ctx, "rtree://")
```

### Database URIs

The `go-whosonfirst-spatial-rtree` package is instantiated using a URI in the form of:

```
rtree://?{PARAMETERS}
```

Valid parameters include:

#### Parameters

| Name | Value | Required| Notes |
| --- | --- | --- | --- |
| strict | bool | N | |
| index_alt_files | bool | N | |

## Tools

```
$> make cli
```

### query

```
$> ./bin/query -h
  -alternate-geometry value
    	One or more alternate geometry labels (wof:alt_label) values to filter results by.
  -custom-placetypes string
    	...
  -custom-placetypes-source string
    	...
  -enable-custom-placetypes
    	...
  -enable-properties
    	Enable support for 'properties' parameters in queries.
  -exclude value
    	Exclude (WOF) records based on their existential flags. Valid options are: ceased, deprecated, not-current, superseded.
  -geometries string
    	Valid options are: all, alt, default. (default "all")
  -index-properties
    	Index properties reader.
  -is-ceased value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-current value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-deprecated value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-superseded value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-superseding value
    	One or more existential flags (-1, 0, 1) to filter results by.
  -is-wof
    	Input data is WOF-flavoured GeoJSON. (Pass a value of '0' or 'false' if you need to index non-WOF documents. (default true)
  -iterator-uri string
    	A valid whosonfirst/go-whosonfirst-iterate/emitter URI. Supported schemes are: directory://, featurecollection://, file://, filelist://, geojsonl://, repo://. (default "repo://")
  -latitude float
    	A valid latitude.
  -longitude float
    	A valid longitude.
  -placetype value
    	One or more place types to filter results by.
  -properties value
    	One or more Who's On First properties to append to each result.
  -properties-reader-uri string
    	Valid options are: []
  -setenv
    	Set flags from environment variables.
  -spatial-database-uri string
    	Valid options are: [rtree://] (default "rtree://")
  -verbose
    	Be chatty.
```

For example:

```
$> ./bin/query \
	-iterator-uri 'repo://?include=properties.mz:is_current=1' \
	-latitude 37.613490350845794 \
	-longitude -122.38882533303682 \
	/usr/local/data/sfomuseum-data-architecture/

| jq '.places[]["wof:name"]'

"Boarding Area A"
"SFO Terminal Complex"
"International Terminal"
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-iterate
* https://github.com/dhconnelly/rtreego