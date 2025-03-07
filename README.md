# go-whosonfirst-spatial-www

Opinionated web application for the `go-whosonfirst-spatial` packages.

## Documentation

Documentation is incomplete at this time.

## Example

```
package main

import (
	_ "github.com/whosonfirst/go-reader-cachereader"
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

There are also server implementations for the following (spatial) databases:

* https://github.com/whosonfirst/go-whosonfirst-spatial-www-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-protomaps
* https://github.com/whosonfirst/go-whosonfirst-spatial-www-duckdb

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

Documentation for the `server` tool has been moved in to [cmd/server/README.md](cmd/server/README.md).

## See also

* https://github.com/whosonfirst/go-whosonfirst-spatial
* https://github.com/whosonfirst/go-whosonfirst-spatial-rtree
* https://github.com/whosonfirst/go-whosonfirst-spatial-pip
* https://github.com/aaronland/go-http-maps