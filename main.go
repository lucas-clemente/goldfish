package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	assetfs "github.com/elazarl/go-bindata-assetfs"

	_ "net/http/pprof"

	"github.com/lucas-clemente/go-http-logger"
	"github.com/lucas-clemente/goldfish/git"
	"github.com/lucas-clemente/goldfish/server"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 3456, "tcp port to listen on")

	flag.Parse()

	if flag.NArg() != 1 {
		log.Fatal("Usage: ./goldfish <path/to/repo>")
	}

	path := flag.Arg(0)

	if err := os.MkdirAll(path, os.ModeDir|0755); err != nil {
		log.Fatal(err)
	}

	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		log.Fatal(err)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := git.NewGitRepo(path)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("🐟 Goldfish listening on http://localhost:%d\n", port)

	http.Handle("/v2/", server.NewHandler2(repo))
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
