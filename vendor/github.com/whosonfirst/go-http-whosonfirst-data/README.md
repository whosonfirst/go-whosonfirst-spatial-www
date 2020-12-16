# go-http-whosonfirst-data

Go HTTP handler for serving Who's On First records.

## Example

```
package main

import (
	"context"
	"github.com/whosonfirst/go-http-whosonfirst-data"
	"github.com/whosonfirst/go-reader"
	"net/http"
)

func main() {

	reader_uri := "fs:///usr/local/data/sfomuseum-data-architecture/data"
	server_uri := "localhost:8080"

	flag.Parse()

	ctx := context.Background()

	r, _ := reader.NewReader(ctx, reader_uri)

	data_handler := data.WhosOnFirstDataHandler(r)

	mux := http.NewServeMux()
	mux.Handle("/", data_handler)

	http.ListenAndServe(server_uri, mux)
}
```

_Error handling removed for the sake of brevity._

## Readers

This package uses the [go-reader.Reader](https://github.com/whosonfirst/go-reader) interface for retrieving data associated with a Who's On First record. The only reader that is available by default, with this package, is the [local filesystem reader](https://github.com/whosonfirst/go-reader#fs) (`fs://`). If you need to use other `go-reader.Reader` implementations in your code you will need to import them explicitly.

## URIs

The `data.WhosOnFirstDataHandler` handler will honour any URL path that can be parsed by the `go-whosonfirst-uri` package. For example :

* localhost:8080/1159157333
* localhost:8080/1159157333.geojson
* localhost:8080/115/915/733/3/1159157333
* localhost:8080/1159157333-alt-sfogis-ground
* localhost:8080/1159157333-alt-sfogis-ground.geojson
* localhost:8080/115/915/733/3/1159157333-alt-sfogis-ground.geojson

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

A simple HTTP server for Who's On First records. The `server` tool is provided as an example application and not necessarily suited for production use.

```
$> ./bin/server -h
Usage of ./bin/server:
  -reader-uri string
    	A valid whosonfirst/go-reader URI
  -server-uri string
    	The host address and port to listen for requests on (default "localhost:8080")
```

For example:

```
$> ./bin/server -reader-uri fs:///usr/local/data/sfomuseum-data-architecture/data
2020/12/16 11:11:49 Listening for requests on localhost:8080
```

And then, from another terminal:

```
$> curl -s localhost:8080/1159157333 | jq '.["properties"]["wof:name"]'
"International Terminal"
```

## See also

* https://github.com/whosonfirst/go-reader
* https://github.com/whosonfirst/go-whosonfirst-uri