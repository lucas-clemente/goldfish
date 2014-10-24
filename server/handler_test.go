package server_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"sort"

	"github.com/lucas-clemente/goldfish/server"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockRepo struct {
	files map[string]string
}

func (r *mockRepo) ReadFile(path string) (io.ReadCloser, error) {
	if c, ok := r.files[path]; ok {
		return ioutil.NopCloser(bytes.NewBufferString(c)), nil
	}
	return nil, os.ErrNotExist
}
func (r *mockRepo) StoreFile(path string, reader io.Reader) error {
	c, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	r.files[path] = string(c)
	return nil
}

func (r *mockRepo) ListFiles(prefix string) ([]string, error) {
	paths := []string{}
	pathRegex := regexp.MustCompile(prefix + "[^/]*/?")
	for p := range r.files {
		if f := pathRegex.FindString(p); f != "" {
			paths = append(paths, f)
		}
	}
	sort.Strings(paths)
	return paths, nil
}

func (r *mockRepo) Observer() <-chan string {
	panic("not implemented")
}

func (r *mockRepo) CloseObserver(<-chan string) {
	panic("not implemented")
}

var _ = Describe("Handler", func() {
	var (
		repo *mockRepo
		resp *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		files := map[string]string{
			"/foo/bar.md": "foobar",
			"/baz":        "foobaz",
		}
		repo = &mockRepo{files: files}
		resp = httptest.NewRecorder()
	})

	It("GETs pages", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/foo/bar.md", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Body.String()).To(Equal("foobar"))
	})

	It("404s", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/noooooooo", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNotFound))
	})

	It("POST updates pages", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("POST", "/v1/foo/bar.md", bytes.NewBufferString("foobaz"))
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNoContent))
		Expect(repo.files["/foo/bar.md"]).To(Equal("foobaz"))
	})

	It("POSTs new pages", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("POST", "/v1/new", bytes.NewBufferString("foobaz"))
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusNoContent))
		Expect(repo.files["/new"]).To(Equal("foobaz"))
	})

	It("GETs root", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))
		Expect(resp.Body.String()).To(MatchJSON(`["/baz", "/foo/"]`))
	})

	It("GETs subdir", func() {
		handler := server.NewHandler(repo, "/v1")
		req, err := http.NewRequest("GET", "/v1/foo/", nil)
		Expect(err).To(BeNil())
		handler.ServeHTTP(resp, req)
		Expect(resp.Code).To(Equal(http.StatusOK))
		Expect(resp.Header().Get("Content-Type")).To(Equal("application/json"))
		Expect(resp.Body.String()).To(MatchJSON(`["/foo/bar.md"]`))
	})
})
