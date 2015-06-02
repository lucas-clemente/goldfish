package server

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lucas-clemente/goldfish/repository"
	"github.com/skratchdot/open-golang/open"
)

// NewHandler2 makes a http.Handler for a given repo.
func NewHandler2(repo repository.Repo) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.GET("/v2/raw/*path", errorHandler(func(c *gin.Context) error {
		path := c.Params.ByName("path")

		f, err := repo.ReadFile(path)
		if err != nil {
			return err
		}

		reader, err := f.Reader()
		if err != nil {
			return err
		}
		defer reader.Close()

		c.Writer.Header().Set("Content-Type", getContentType(path))

		if _, err := io.Copy(c.Writer, reader); err != nil {
			log.Println(err)
		}
		return nil
	}))

	router.POST("/v2/raw/*path", errorHandler(func(c *gin.Context) error {
		path := c.Params.ByName("path")

		err := repo.StoreFile(path, c.Request.Body)
		if err != nil {
			return err
		}

		c.Writer.WriteHeader(http.StatusNoContent)
		return nil
	}))

	router.GET("/v2/folders/:id", errorHandler(func(c *gin.Context) error {
		id := c.Params.ByName("id")

		entries, err := repo.ListFiles(idToPath(id))
		if err != nil {
			return err
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

		c.JSON(200, map[string]interface{}{
			"folder": map[string]interface{}{
				"id":           id,
				"pages":        pageIDs,
				"subfolders":   subfolderIDs,
				"parentFolder": parentID,
			},
		})
		return nil
	}))

	router.GET("/v2/pages/:id", errorHandler(func(c *gin.Context) error {
		id := c.Params.ByName("id")

		file, err := repo.ReadFile(idToPath(id))
		if err != nil {
			return err
		}

		jsonPage, err := getPageJSON(file)
		if err != nil {
			return err
		}

		c.JSON(200, map[string]interface{}{"page": jsonPage})
		return nil
	}))

	router.DELETE("/v2/pages/:id", errorHandler(func(c *gin.Context) error {
		id := c.Params.ByName("id")

		if err := repo.DeleteFile(idToPath(id)); err != nil {
			return err
		}
		c.Writer.WriteHeader(http.StatusNoContent)
		return nil
	}))

	router.POST("/v2/pages/:id/open", errorHandler(func(c *gin.Context) error {
		id := c.Params.ByName("id")

		p, err := repo.LocalPathForFile(idToPath(id))

		if err != nil {
			return err
		}

		if err := open.Run(p); err != nil {
			return err
		}

		c.Writer.WriteHeader(http.StatusNoContent)
		return nil
	}))

	router.GET("/v2/pages", errorHandler(func(c *gin.Context) error {
		var searchTerm string
		searchTermList, ok := c.Request.URL.Query()["q"]
		if ok && len(searchTermList) != 0 {
			searchTerm = searchTermList[0]
		}

		results, err := repo.SearchFiles(searchTerm)
		if err != nil {
			return err
		}

		jsonArray := []interface{}{}
		for _, file := range results {
			jsonPage, err := getPageJSON(file)
			if err != nil {
				return err
			}
			jsonArray = append(jsonArray, jsonPage)
		}

		c.JSON(200, map[string]interface{}{"pages": jsonArray})
		return nil
	}))

	return router
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
	case "pdf":
		return "application/pdf"
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

func errorHandler(h func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			if os.IsNotExist(err) {
				c.AbortWithError(404, err)
			}
			log.Println(err)
			c.AbortWithError(500, err)
		}
	}
}
