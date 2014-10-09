all: build

deps:
	go get -t

build:
	GOOS=darwin go build -o build/notes_darwin
	GOOS=windows go build -o build/notes.exe
	GOOS=linux go build -o build/notes_linux
