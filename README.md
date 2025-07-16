# go-whosonfirst-spatial-www

Opinionated web application for the `go-whosonfirst-spatial` packages.

## Documentation

Documentation is incomplete at this time.

## Example

```
package main

import (
	"context"
	"log"

	_ "github.com/whosonfirst/go-reader-cachereader"
	"github.com/whosonfirst/go-whosonfirst-spatial-www/application/server"
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

* https://github.com/whosonfirst/go-whosonfirst-spatial-sqlite
* https://github.com/whosonfirst/go-whosonfirst-spatial-pmtiles
* https://github.com/whosonfirst/go-whosonfirst-spatial-duckdb

## Tools

```
$> make cli
go build -mod vendor -o bin/server cmd/server/main.go
```

### server

Documentation for the `server` tool has been moved in to [cmd/server/README.md](cmd/server/README.md).
