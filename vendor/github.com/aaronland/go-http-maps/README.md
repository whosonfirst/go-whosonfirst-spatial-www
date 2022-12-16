# go-http-maps

Go package providing opinionated HTTP middleware for web-based map tiles.

## Important

This is work in progress. Documentation is incomplete.

Until then have a look at [app/server/main.go](app/server/main.go), [templates/html/map.html](templates/html/map.html) and [static/javascript/aaronland.map.init.js](static/javascript/aaronland.map.init.js) for an example of working code.

## Example

![](docs/images/go-http-maps-radius.png)

```
go run -mod vendor cmd/server/main.go \
	-map-provider tangramjs \
	-tilezen-enable-tilepack \
	-tilezen-tilepack-path /usr/local/data/sf.db
```

## See also

* https://github.com/aaronland/go-http-leaflet
* https://github.com/aaronland/go-http-tangramjs
* https://github.com/aaronland/go-http-protomaps
* https://github.com/tilezen/go-tilepacks
