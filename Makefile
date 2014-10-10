.PHONY: all deps debug web bindata build

all: build

deps:
	go get -t

debug:
	go-bindata -prefix=web/dist -debug=true web/dist/...

web:
	cd web && ember build --production

bindata: web
	go-bindata -prefix=web/dist web/dist/...

build: bindata
	go build -o build/notes
