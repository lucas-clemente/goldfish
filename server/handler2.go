package server

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// NewHandler2 makes a http.Handler for a given repo.
func NewHandler2(repo Repo) http.Handler {
	router := httprouter.New()

	router.GET("/v2/raw/*path", func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		path := p.ByName("path")

		c, err := repo.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, nil)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		defer c.Close()
		if _, err := io.Copy(w, c); err != nil {
			log.Println(err)
		}
	})

	router.GET("/v2/folders/*id", func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		id := strings.TrimLeft(p.ByName("id"), "/")

		entries, err := repo.ListFiles(idToPath(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		pageIDs := []string{}
		subfolderIDs := []string{}

		for _, entry := range entries {
			if entry[len(entry)-1] == '/' {
				subfolderIDs = append(subfolderIDs, pathToID(entry[0:len(entry)-1]))
			} else {
				pageIDs = append(pageIDs, pathToID(entry))
			}
		}

		var parentID interface{}
		if id != "|" {
			parentID = id[0:strings.LastIndex(id, "|")]
			if parentID == "" {
				parentID = "|"
			}
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"folder": map[string]interface{}{
				"id":           id,
				"pages":        pageIDs,
				"subfolders":   subfolderIDs,
				"parentFolder": parentID,
			},
		})
		if err != nil {
			log.Println(err)
		}
	})

	router.GET("/v2/pages/*id", func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		id := strings.TrimLeft(p.ByName("id"), "/")

		c, err := repo.ReadFile(idToPath(id))
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, nil)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		c.Close()

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"page": map[string]interface{}{
				"id": id,
			},
		})
		if err != nil {
			log.Println(err)
		}
	})

	return router
}

func handleError2(err error, w http.ResponseWriter) {
	if os.IsNotExist(err) {
		http.NotFound(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func idToPath(id string) string {
	return strings.Replace(id, "|", "/", -1)
}

func pathToID(path string) string {
	return strings.Replace(path, "/", "|", -1)
}
