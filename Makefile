cli:
	go build -mod vendor -o bin/server cmd/server/main.go

debug:
	go run -mod vendor cmd/server/main.go -enable-www -enable-tangram -spatial-database-uri 'rtree:///?strict=false&index_alt_files=1' -properties-reader-uri 'cachereader://?reader=repo:///$(REPO)&cache=gocache://'  -nextzen-apikey $(APIKEY) -iterator-uri repo:// $(REPO)
