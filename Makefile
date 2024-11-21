GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

vuln:
	govulncheck ./...

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/server cmd/server/main.go

debug:
	go run -mod $(GOMOD) cmd/server/main.go \
		-enable-www \
		-spatial-database-uri 'rtree:///?strict=false&index_alt_files=0' \
		-properties-reader-uri 'cachereader://?reader=repo://$(REPO)&cache=gocache://' \
		-iterator-uri 'repo://#$(REPO)'
