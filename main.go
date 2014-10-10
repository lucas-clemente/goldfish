package main

import (
	"log"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/lucas-clemente/notes/notes"
)

const path = "tmp/repo"

func main() {
	repo, err := notes.NewGitRepo(path)
	if err != nil {
		log.Fatal(err)
	}
	_ = repo

	http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir}))
	log.Fatal(http.ListenAndServe("localhost:3456", nil))
}
