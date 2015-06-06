.PHONY: all debug web bindata build build-release ci-install test

all: build

debug:
	mkdir -p web/dist
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

ci-install:
	# Git
	git config --global user.name "Test User"
	git config --global user.email "test@example.com"
	# Go build
	go get github.com/jteeuwen/go-bindata/...
	$(MAKE) debug
	go get -t ./...
	# JS
	npm install -g ember-cli bower
	cd web && npm install
	cd web && bower install

test:
	go test -v ./...
	cd web && ember test
