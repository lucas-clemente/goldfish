package git

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	git2go "github.com/libgit2/git2go"
)

// GitRepo is a git repository implementing the Repo interface for goldfish.
type GitRepo struct {
	path string
	repo *git2go.Repository
}

// NewGitRepo opens or makes a git repo at the given path
func NewGitRepo(path string) (*GitRepo, error) {
	repo, err := git2go.OpenRepository(path)
	if err != nil {
		repo, err = git2go.InitRepository(path, false)
		if err != nil {
			return nil, err
		}

		// Make empty tree
		index, err := repo.Index()
		if err != nil {
			return nil, err
		}
		defer index.Free()

		treeID, err := index.WriteTree()
		if err != nil {
			return nil, err
		}

		tree, err := repo.LookupTree(treeID)
		if err != nil {
			return nil, err
		}

		defer tree.Free()
		sig := &git2go.Signature{Name: "system", Email: "notes@clemente.io", When: time.Now()}
		_, err = repo.CreateCommit("refs/heads/master", sig, sig, "initial commit", tree)
		if err != nil {
			return nil, err
		}
	}

	return &GitRepo{path: path, repo: repo}, nil
}

// ReadFile reads a file from the repo
func (r *GitRepo) ReadFile(path string) (io.ReadCloser, error) {
	return os.Open(r.absolutePath(path))
}

// StoreFile writes a file to the repo and commits it
func (r *GitRepo) StoreFile(p string, data io.Reader) error {
	if err := os.MkdirAll(path.Dir(r.absolutePath(p)), 0755); err != nil {
		return err
	}

	file, err := os.OpenFile(r.absolutePath(p), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return err
	}

	index, err := r.repo.Index()
	if err != nil {
		return err
	}
	defer index.Free()

	if err := index.AddByPath(p); err != nil {
		return err
	}

	treeID, err := index.WriteTree()
	if err != nil {
		return err
	}

	return r.commit(treeID, p)
}

// ListFiles lists the files in a given directory
func (r *GitRepo) ListFiles(prefix string) ([]string, error) {
	files := []string{}

	commit, err := r.headCommit()
	if err != nil {
		return nil, err
	}
	defer commit.Free()

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}
	defer tree.Free()

	if prefix != "/" {
		prefixTreeID, err := tree.EntryByPath(strings.TrimPrefix(prefix, "/"))
		if err != nil {
			return nil, err
		}

		tree, err = r.repo.LookupTree(prefixTreeID.Id)
		if err != nil {
			return nil, err
		}
		defer tree.Free()
	}

	var i uint64
	for i = 0; i < tree.EntryCount(); i++ {
		entry := tree.EntryByIndex(i)
		f := prefix + entry.Name
		switch entry.Type {
		case git2go.ObjectBlob:
			files = append(files, f)
		case git2go.ObjectTree:
			files = append(files, f+"/")
		default:
			panic("unexpected object in tree")
		}
	}

	return files, nil
}

func (r *GitRepo) headCommit() (*git2go.Commit, error) {
	headRef, err := r.repo.Head()
	defer headRef.Free()
	if err != nil {
		log.Fatal(err)
	}

	headID := headRef.Target()
	return r.repo.LookupCommit(headID)
}

func (r *GitRepo) commit(treeID *git2go.Oid, message string) error {
	tree, err := r.repo.LookupTree(treeID)
	if err != nil {
		return err
	}
	defer tree.Free()

	headCommit, err := r.headCommit()
	if err != nil {
		return err
	}
	defer headCommit.Free()

	sig := &git2go.Signature{Name: "system", Email: "notes@clemente.io", When: time.Now()}
	_, err = r.repo.CreateCommit("refs/heads/master", sig, sig, message, tree, headCommit)
	if err != nil {
		return err
	}

	return nil
}

func (r *GitRepo) absolutePath(path string) string {
	return r.path + path
}
