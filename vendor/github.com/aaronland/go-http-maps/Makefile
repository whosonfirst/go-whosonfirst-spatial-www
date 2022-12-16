CWD=$(shell pwd)

cli:
	go build -o bin/server cmd/server/main.go

debug-tangram:
	go run -mod vendor cmd/server/main.go -map-provider tangram -nextzen-apikey $(APIKEY)

debug-tilepack:
	go run -mod vendor cmd/server/main.go -map-provider tangram -tilezen-enable-tilepack -tilezen-tilepack-path /usr/local/data/sf.db

debug-protomaps:
	go run -mod vendor cmd/server/main.go -map-provider protomaps -protomaps-serve-tiles -protomaps-bucket-uri file://$(CWD)/fixtures -protomaps-database sfo

debug-leaflet:
	go run -mod vendor cmd/server/main.go -map-provider leaflet -leaflet-tile-url https://tile.openstreetmap.org/{z}/{x}/{y}.png
