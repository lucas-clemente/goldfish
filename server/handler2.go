package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

		w.Header().Set("Content-Type", getContentType(path))

		if _, err := io.Copy(w, c); err != nil {
			log.Println(err)
		}
	})

	router.POST("/v2/raw/*path", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		path := p.ByName("path")

		err := repo.StoreFile(path, r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
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
		defer c.Close()

		folder := id[0:strings.LastIndex(id, "|")]
		if folder == "" {
			folder = "|"
		}

		var markdownSource interface{}
		if strings.HasSuffix(id, ".md") {
			markdownSourceBytes, err := ioutil.ReadAll(c)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			markdownSource = string(markdownSourceBytes)
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"page": map[string]interface{}{
				"id":             id,
				"folder":         folder,
				"markdownSource": markdownSource,
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

func getContentType(filename string) string {
	extension := filename[strings.LastIndex(filename, ".")+1:]
	switch extension {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "svg":
		return "image/svg+xml"
	}
	return "text/plain"
}
