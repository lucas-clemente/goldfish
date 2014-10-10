package notes

// Repo is a collection of files that make up a wiki
type Repo interface {
	ReadFile(path string) ([]byte, error)
	StoreFile(path string, data []byte) error
}
