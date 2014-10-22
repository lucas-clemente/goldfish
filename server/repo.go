package server

import "io"

// Repo is a collection of files that make up a wiki
type Repo interface {
	ReadFile(path string) (io.ReadCloser, error)
	StoreFile(path string, content io.Reader) error
	ListFiles(prefix string) ([]string, error)
}
