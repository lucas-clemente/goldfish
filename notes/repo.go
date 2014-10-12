package notes

import "io"

// NotFoundError can be returned by Repo
type NotFoundError struct{}

func (NotFoundError) Error() string {
	return "not found"
}

// Repo is a collection of files that make up a wiki
type Repo interface {
	ReadFile(path string) (io.ReadCloser, error)
	StoreFile(path string, content io.Reader) error
}
