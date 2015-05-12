package server_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"

	"github.com/lucas-clemente/goldfish/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockRepo2 struct {
	files map[string]string
}

func (r *mockRepo2) ReadFile(path string) (io.ReadCloser, error) {
	if c, ok := r.files[path]; ok {
		return ioutil.NopCloser(bytes.NewBufferString(c)), nil
	}
	return nil, os.ErrNotExist
}
func (r *mockRepo2) StoreFile(path string, reader io.Reader) error {
	c, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	r.files[path] = string(c)
	return nil
}

func (r *mockRepo2) ListFiles(prefix string) ([]string, error) {
	paths := []string{}
	for p := range r.files {
		if strings.HasPrefix(p, prefix) {
			withoutPrefix := strings.TrimPrefix(p, prefix+"/")
			if strings.Contains(withoutPrefix, "/") {
				paths = append(paths, prefix+"/"+strings.Split(withoutPrefix, "/")[0]+"/")
			} else {
				paths = append(paths, p)
			}
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (r *mockRepo2) Observer() <-chan string {
	panic("not implemented")
}

func (r *mockRepo2) CloseObserver(<-chan string) {
	panic("not implemented")
}

var _ = Describe("Handler", func() {
	var (
		repo *mockRepo2
		resp *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		files := map[string]string{
			"/foo/bar.md":     "foobar",
			"/foo/fuu/bar.md": "foobar",
			"/baz":            "foobaz",
		}
		repo = &mockRepo2{files: files}
		resp = httptest.NewRecorder()
	})

	It("GETs raw", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/raw/foo/bar.md", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(Equal("foobar"))
	})

	It("404s", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/noooooooo", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNotFound))
	})

	It("GETs folders", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/folders/foo", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"folder":{"id":"/foo","pages":["/foo/bar.md"],"subfolders":["/foo/fuu"]}}`))
	})
})
