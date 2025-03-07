GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

REPO=/usr/local/data/sfomuseum-data-architecture
INITIAL_VIEW=-122.384292,37.621131,13

vuln:
	govulncheck ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-enable-www \
		-enable-geojson \
		-spatial-database-uri 'rtree:///?strict=false&index_alt_files=0' \
		-properties-reader-uri 'cachereader://?reader=repo://$(REPO)&cache=gocache://' \
		-initial-view '$(INITIAL_VIEW)' \
		-iterator-uri 'repo://#$(REPO)'
