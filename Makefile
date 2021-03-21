cli:
	go build -mod vendor -o bin/server cmd/server/main.go

debug:
	go run -mod vendor cmd/server/main.go -enable-www -index-properties -spatial-database-uri 'rtree:///?strict=false&index_alt_files=1' -properties-reader-uri 'whosonfirst://?reader=fs:///$(REPO)/data&cache=gocache://' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-woeplanet:
	go run -mod vendor cmd/server/main.go -enable-www -properties-reader-uri 'whosonfirst://?reader=fs:///$(REPO)/data&cache=gocache://' -geojson-path-resolver-uri wofid:// -spatial-database-uri 'rtree:///?strict=false&index_alt_files=1' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-geojson:
	go run -mod vendor cmd/server/main.go -spatial-database-uri 'rtree:///?strict=false' -mode featurecollection:// $(GEOJSON)

