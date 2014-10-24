package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"code.google.com/p/go.net/websocket"
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
	if len(path) == 0 {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET":
		if path[len(path)-1] == '/' {
			// Folder
			files, err := h.repo.ListFiles(path)
			if err != nil {
				handleError(err, w)
				return
			}

			data, err := json.Marshal(files)
			if err != nil {
				handleError(err, w)
				return
			}

			w.Header().Set("Content-Type", "application/json")

			w.Write(data)
		} else if r.Header.Get("Connection") == "Upgrade" {
			// Page via websocket
			websocket.Handler(func(conn *websocket.Conn) {
				changes := h.repo.Observer()
				defer h.repo.CloseObserver(changes)
				if err := h.serveFile(path, conn); err != nil {
					log.Printf("error in ws: %s\n", err.Error())
				}
				for changedFile := range changes {
					if changedFile == path {
						if err := h.serveFile(path, conn); err != nil {
							log.Printf("error in ws: %s\n", err.Error())
						}
					}
				}
			}).ServeHTTP(w, r)
		} else {
			// Normal file
			if err := h.serveFile(path, w); err != nil {
				handleError(err, w)
			}
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

func handleError(err error, w http.ResponseWriter) {
	if os.IsNotExist(err) {
		http.NotFound(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *handler) serveFile(path string, w io.Writer) error {
	c, err := h.repo.ReadFile(path)
	if err != nil {
		return err
	}
	defer c.Close()
	if _, err := io.Copy(w, c); err != nil {
		log.Println(err)
	}
	return nil
}
