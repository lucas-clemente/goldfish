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
	cd web && ember build --environment=production

bindata: web
	go-bindata -nomemcopy=true -prefix=web/dist web/dist/...

build: bindata
	go build -o build/goldfish

build-release: bindata
	GOOS=darwin go build -o build/goldfish_darwin
	cd build && cp goldfish_darwin goldfish && zip goldfish.darwin.zip goldfish
	GOOS=windows go build -o build/goldfish_windows
	cd build && cp goldfish_windows goldfish.exe && zip goldfish.windows.zip goldfish.exe
	GOOS=linux go build -o build/goldfish_linux
	cd build && cp goldfish_linux goldfish && zip goldfish.linux.zip goldfish
