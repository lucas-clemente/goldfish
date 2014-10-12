package notes

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type handler struct {
	repo   Repo
	prefix string
}

// NewHandler makes a http.Handler for a given repo.
func NewHandler(repo Repo, prefix string) http.Handler {
	return &handler{repo: repo, prefix: prefix}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, h.prefix)

	switch r.Method {
	case "GET":
		c, err := h.repo.ReadFile(path)
		if err != nil {
			if _, ok := err.(NotFoundError); ok {
				http.NotFound(w, r)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		defer c.Close()
		if _, err := io.Copy(w, c); err != nil {
			log.Println(err)
		}

	case "POST":
		defer r.Body.Close()
		err := h.repo.StoreFile(path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
