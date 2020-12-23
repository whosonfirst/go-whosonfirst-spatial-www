# go-whosonfirst-spatial-rtree

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
  -latitude float
    	A valid latitude.
  -longitude float
    	A valid longitude.
  -mode string
    	Valid modes are: directory, featurecollection, file, filelist, geojsonl, repo. (default "repo://")
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
	-database-uri 'rtree://?strict=false' \
	-latitude 37.616951 \
	-longitude -122.383747 \
	-mode repo:// \
	/usr/local/data/sfomuseum-data-architecture/ \

| jq | grep wof:name

17:08:24.974105 [query][index] ERROR 1159157931 failed indexing, (rtreego: improper distance). Strict mode is disabled, so skipping.
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "International Terminal",
      "wof:name": "International Terminal",
      "wof:name": "Central Terminal",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "Central Terminal",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "Terminal 2 Main Hall",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "Central Terminal",
      "wof:name": "Terminal 2",
      "wof:name": "Terminal 2 Main Hall",
      "wof:name": "Terminal 2",
      "wof:name": "Central Terminal",
      "wof:name": "Boarding Area D",
      "wof:name": "Boarding Area D",
      "wof:name": "Central Terminal",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
      "wof:name": "SFO Terminal Complex",
```

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/dhconnelly/rtreego