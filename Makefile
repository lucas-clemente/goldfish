.PHONY: all deps debug web bindata build

all: build

deps:
	go get -t
	go get github.com/jteeuwen/go-bindata/...
	cd web && npm install
	cd web && bower install

debug:
	go-bindata -prefix=web/dist -debug=true web/dist/...

web:
	rm -rf web/dist
	cd web && ENV=production gulp

bindata: web
	go-bindata -prefix=web/dist web/dist/...

build: bindata
	go build -o build/notes
