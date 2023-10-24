GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

vuln:
	govulncheck ./...

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/server cmd/server/main.go

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-enable-www \
		-map-provider 'leaflet://' \
		-spatial-database-uri 'rtree:///?strict=false&index_alt_files=0' \
		-properties-reader-uri 'cachereader://?reader=repo://$(REPO)&cache=gocache://' \
		-nextzen-apikey $(APIKEY) \
		-iterator-uri 'repo://' \
		$(REPO)
