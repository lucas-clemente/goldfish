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
	"github.com/lucas-clemente/goldfish/repository"
)

// NewHandler2 makes a http.Handler for a given repo.
func NewHandler2(repo repository.Repo) http.Handler {
	router := httprouter.New()

	router.GET("/v2/raw/*path", func(w http.ResponseWriter, _ *http.Request, p httprouter.Params) {
		path := p.ByName("path")

		f, err := repo.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, nil)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		reader, err := f.Reader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer reader.Close()

		w.Header().Set("Content-Type", getContentType(path))

		if _, err := io.Copy(w, reader); err != nil {
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
			entryPath := entry.Path()
			if entryPath[len(entryPath)-1] == '/' {
				subfolderIDs = append(subfolderIDs, pathToID(entryPath[0:len(entryPath)-1]))
			} else {
				pageIDs = append(pageIDs, pathToID(entryPath))
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

		file, err := repo.ReadFile(idToPath(id))
		if err != nil {
			if os.IsNotExist(err) {
				http.NotFound(w, nil)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		jsonPage, err := getPageJSON(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{"page": jsonPage})
		if err != nil {
			log.Println(err)
		}
	})

	router.GET("/v2/pages", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var searchTerm string
		searchTermList, ok := r.URL.Query()["q"]
		if ok && len(searchTermList) != 0 {
			searchTerm = searchTermList[0]
		}

		results, err := repo.SearchFiles(searchTerm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonArray := []interface{}{}
		for _, file := range results {
			jsonPage, err := getPageJSON(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonArray = append(jsonArray, jsonPage)
		}

		err = json.NewEncoder(w).Encode(map[string]interface{}{"pages": jsonArray})
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

func getPageJSON(file repository.File) (interface{}, error) {
	id := pathToID(file.Path())

	folder := id[0:strings.LastIndex(id, "|")]
	if folder == "" {
		folder = "|"
	}

	var markdownSource interface{}
	if strings.HasSuffix(id, ".md") {
		c, err := file.Reader()
		if err != nil {
			return nil, err
		}
		defer c.Close()

		markdownSourceBytes, err := ioutil.ReadAll(c)
		if err != nil {
			return nil, err
		}
		markdownSource = string(markdownSourceBytes)
	}

	return (map[string]interface{}{
		"id":             id,
		"folder":         folder,
		"modifiedAt":     file.ModTime(),
		"markdownSource": markdownSource,
	}), nil
}
