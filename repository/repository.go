package repository

import (
	"io"
	"time"
)

// A File in a repo
type File interface {
	Path() string
	Reader() (io.ReadCloser, error)
	ModTime() time.Time
}

// Repo is a collection of files that make up a wiki
type Repo interface {
	ReadFile(path string) (File, error)
	StoreFile(path string, content io.Reader) error
	DeleteFile(path string) error
	ListFiles(prefix string) ([]File, error)
	SearchFiles(term string) ([]File, error)

	Observer() <-chan string
	CloseObserver(<-chan string)
}
