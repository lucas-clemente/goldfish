package server_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"github.com/lucas-clemente/goldfish/repository"
	"github.com/lucas-clemente/goldfish/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockFile struct {
	path     string
	modTime  time.Time
	contents string
}

var _ repository.File = &mockFile{}

func (f *mockFile) Path() string {
	return f.path
}

func (f *mockFile) Reader() (io.ReadCloser, error) {
	return ioutil.NopCloser(bytes.NewBufferString(f.contents)), nil
}

func (f *mockFile) ModTime() time.Time {
	return f.modTime
}

type mockRepo2 struct {
	files map[string]*mockFile
}

var _ repository.Repo = &mockRepo2{}

func (r *mockRepo2) ReadFile(path string) (repository.File, error) {
	if f, ok := r.files[path]; ok {
		return f, nil
	}
	return nil, os.ErrNotExist
}

func (r *mockRepo2) StoreFile(path string, reader io.Reader) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	r.files[path] = &mockFile{path: path, contents: string(data)}
	return nil
}

func (r *mockRepo2) ListFiles(prefix string) ([]repository.File, error) {
	if prefix == "/" {
		return []repository.File{r.files["/baz"], &mockFile{path: "/foo/"}}, nil
	} else if prefix == "/foo" {
		return []repository.File{r.files["/foo/bar.md"], &mockFile{path: "/foo/fuu/"}}, nil
	} else if prefix == "/foo/fuu" {
		return []repository.File{r.files["/foo/fuu/bar.md"]}, nil
	}
	return nil, os.ErrNotExist
}

func (r *mockRepo2) SearchFiles(term string) ([]repository.File, error) {
	if term != "foobar" {
		return []repository.File{}, nil
	}
	return []repository.File{r.files["/foo/bar.md"], r.files["/foo/fuu/bar.md"]}, nil
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
			files: map[string]*mockFile{
				"/foo/bar.md":     &mockFile{path: "/foo/bar.md", contents: "foobar"},
				"/foo/fuu/bar.md": &mockFile{path: "/foo/fuu/bar.md", contents: "foobar"},
				"/baz":            &mockFile{path: "/baz", contents: "foobaz"},
			},
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
		Expect(repo.files["/baz"].contents).To(Equal("new content"))
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
		Expect(resp.Body.String()).To(MatchJSON(`{"page":{"id":"|baz","folder":"|","markdownSource":null,"modifiedAt": "0001-01-01T00:00:00Z"}}`))
	})

	It("GETs markdown pages", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/pages/|foo|bar.md", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"page":{"id":"|foo|bar.md","folder":"|foo","markdownSource":"foobar","modifiedAt": "0001-01-01T00:00:00Z"}}`))
	})

	It("searches markdown pages", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/pages?q=foobar", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"pages":[{"id":"|foo|bar.md","folder":"|foo","markdownSource":"foobar","modifiedAt": "0001-01-01T00:00:00Z"},{"id":"|foo|fuu|bar.md","folder":"|foo|fuu","markdownSource":"foobar","modifiedAt": "0001-01-01T00:00:00Z"}]}`))
	})

	It("searches with no results", func() {
		handler := server.NewHandler2(repo)
		req, err := http.NewRequest("GET", "/v2/pages?q=asdasdasdasdasd", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(MatchJSON(`{"pages":[]}`))
	})
})
