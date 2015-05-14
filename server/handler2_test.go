package server_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/lucas-clemente/goldfish/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockRepo2 struct {
	storedData map[string]string
}

func (r *mockRepo2) ReadFile(path string) (io.ReadCloser, error) {
	files := map[string]string{
		"/foo/bar.md":     "foobar",
		"/foo/fuu/bar.md": "foobar",
		"/baz":            "foobaz",
	}

	if c, ok := files[path]; ok {
		return ioutil.NopCloser(bytes.NewBufferString(c)), nil
	}
	return nil, os.ErrNotExist
}
func (r *mockRepo2) StoreFile(path string, reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	r.storedData[path] = string(data)
	return nil
}

func (r *mockRepo2) ListFiles(prefix string) ([]string, error) {
	if prefix == "/" {
		return []string{"/baz", "/foo/"}, nil
	} else if prefix == "/foo" {
		return []string{"/foo/bar.md", "/foo/fuu/"}, nil
	} else if prefix == "/foo/fuu" {
		return []string{"/foo/fuu/bar.md"}, nil
	}
	return nil, os.ErrNotExist
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
		repo = &mockRepo2{
			storedData: map[string]string{},
		}
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

	It("POSTs raw", func() {
		handler := server.NewHandler2(repo)
		body := bytes.NewBufferString("new content")
		req, err := http.NewRequest("POST", "/v2/raw/baz", body)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNoContent))
		Expect(repo.storedData["/baz"]).To(Equal("new content"))
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
		req, err := http.NewRequest("GET", "/v2/folders/|foo", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"folder":{"id":"|foo","pages":["|foo|bar.md"],"subfolders":["|foo|fuu"],"parentFolder":"|"}}`))
	})

	It("GETs nested folders", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/folders/|foo|fuu", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"folder":{"id":"|foo|fuu","pages":["|foo|fuu|bar.md"],"subfolders":[],"parentFolder":"|foo"}}`))
	})

	It("GETs root folder", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/folders/|", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"folder":{"id":"|","pages":["|baz"],"subfolders":["|foo"],"parentFolder":null}}`))
	})

	It("GETs pages", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/pages/|baz", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"page":{"id":"|baz","folder":"|"}}`))
	})
})
