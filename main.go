package main

import (
	"log"
	"net/http"
	"os"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/lucas-clemente/notes/notes"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./notes <path/to/repo>")
	}

	path := os.Args[1]

	repo, err := notes.NewGitRepo(path)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/v1/", notes.NewHandler(repo, "/v1"))

	http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir}))
	log.Fatal(http.ListenAndServe("localhost:3456", nil))
}
