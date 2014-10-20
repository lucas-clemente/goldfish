package main

import (
	"log"
	"net/http"
	"os"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/lucas-clemente/notes/notes"
)

type logger struct {
	http.ResponseWriter
	req    *http.Request
	status int
	start  time.Time
}

func newLogger(w http.ResponseWriter, req *http.Request) *logger {
	return &logger{w, req, 0, time.Now()}
}

func (w *logger) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *logger) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	return w.ResponseWriter.Write(b)
}

func (w *logger) print() {
	duration := int64(time.Since(w.start) / time.Microsecond)
	log.Printf("[%d] %s %s in %d us\n", w.status, w.req.Method, w.req.RequestURI, duration)
}

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

	log.Fatal(
		http.ListenAndServe("localhost:3456", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			logger := newLogger(w, req)
			http.DefaultServeMux.ServeHTTP(logger, req)
			logger.print()
		})),
	)
}
