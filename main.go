package main

import (
	"github.com/lucas-clemente/notes/notes"
)

const path = "tmp/repo"

func main() {
	repo, err := notes.NewRepo(path)
}
