package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	"github.com/lucas-clemente/go-http-logger"
	"github.com/lucas-clemente/goldfish/git"
	"github.com/lucas-clemente/goldfish/server"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./goldfish <path/to/repo>")
	}

	path := os.Args[1]

	repo, err := git.NewGitRepo(path)
	if err != nil {
		log.Fatal(err)
	}

	port := 3456

	fmt.Printf("🐟 Goldfish listening on http://localhost:%d\n", port)

	http.Handle("/v1/", server.NewHandler(repo, "/v1"))
	http.Handle("/assets/", http.FileServer(&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir}))
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		index, err := Asset("index.html")
		if err != nil {
			panic("could not find index.html")
		}
		w.Write(index)
	})

	log.Fatal(
		http.ListenAndServe("localhost:"+strconv.Itoa(port), logger.Logger(http.DefaultServeMux)),
	)
}
