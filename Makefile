CWD=$(shell pwd)

go-bindata:
	mkdir -p cmd/go-bindata
	mkdir -p cmd/go-bindata-assetfs
	curl -s -o cmd/go-bindata/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata/master/cmd/go-bindata/main.go
	curl -s -o cmd/go-bindata-assetfs/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata-assetfs/master/cmd/go-bindata-assetfs/main.go

bake:
	@make bake-static
	@make bake-templates

bake-static:
	go build -mod vendor -o bin/go-bindata cmd/go-bindata/main.go
	go build -mod vendor -o bin/go-bindata-assetfs cmd/go-bindata-assetfs/main.go
	rm -f static/*~ static/css/*~ static/javascript/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg http -o http/staticfs.go static static/javascript/ static/css/

bake-templates:
	rm -rf templates/html/*~
	bin/go-bindata -pkg templates -o assets/templates/html.go templates/html

cli:
	go build -mod vendor -o bin/server cmd/server/main.go

debug:
	go run -mod vendor cmd/server/main.go -enable-www -index-properties -spatial-database-uri 'rtree:///?strict=false&index_alt_files=1' -properties-reader-uri 'whosonfirst://?reader=fs:///$(REPO)/data&cache=gocache://' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-woeplanet:
	go run -mod vendor cmd/server/main.go -enable-www -properties-reader-uri 'whosonfirst://?reader=fs:///$(REPO)/data&cache=gocache://' -geojson-path-resolver-uri wofid:// -spatial-database-uri 'rtree:///?strict=false&index_alt_files=1' -geojson-reader-uri 'fs://$(REPO)/data' -nextzen-apikey $(APIKEY) -mode repo:// $(REPO)

debug-geojson:
	go run -mod vendor cmd/server/main.go -spatial-database-uri 'rtree:///?strict=false' -mode featurecollection:// $(GEOJSON)

