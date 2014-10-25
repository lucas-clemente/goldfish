.PHONY: all deps debug web bindata build

GOOS = $(shell go env GOOS)

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
	go build -o build/goldfish_$(GOOS)
	cd build && cp goldfish_$(GOOS) goldfish && zip goldfish.$(GOOS).zip goldfish
