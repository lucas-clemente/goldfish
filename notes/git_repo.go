package notes

import (
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"

	git "github.com/libgit2/git2go"
)

type gitRepo struct {
	path string
	repo *git.Repository
}

// NewGitRepo opens or makes a git repo at the given path
func NewGitRepo(path string) (Repo, error) {
	repo, err := git.OpenRepository(path)
	if err != nil {
		repo, err = git.InitRepository(path, false)
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
		sig := &git.Signature{Name: "system", Email: "notes@clemente.io", When: time.Now()}
		_, err = repo.CreateCommit("refs/heads/master", sig, sig, "initial commit", tree)
		if err != nil {
			return nil, err
		}
	}

	return &gitRepo{path: path, repo: repo}, nil
}

func (r *gitRepo) ReadFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(r.absolutePath(path))
	if os.IsNotExist(err) {
		return nil, NotFoundError{}
	}
	return f, err
}

func (r *gitRepo) StoreFile(p string, data io.Reader) error {
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

func (r *gitRepo) ListFiles(prefix string) ([]string, error) {
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

	files := []string{}

	// TODO use prefix to chose initial tree
	err = tree.Walk(func(path string, e *git.TreeEntry) int {
		f := "/" + path + e.Name
		if e.Type == git.ObjectBlob && strings.HasPrefix(f, prefix) {
			files = append(files, f)
		}
		return 0
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (r *gitRepo) headCommit() (*git.Commit, error) {
	headRef, err := r.repo.Head()
	defer headRef.Free()
	if err != nil {
		log.Fatal(err)
	}

	headID := headRef.Target()
	return r.repo.LookupCommit(headID)
}

func (r *gitRepo) commit(treeID *git.Oid, message string) error {
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

	sig := &git.Signature{Name: "system", Email: "notes@clemente.io", When: time.Now()}
	_, err = r.repo.CreateCommit("refs/heads/master", sig, sig, message, tree, headCommit)
	if err != nil {
		return err
	}

	return nil
}

func (r *gitRepo) absolutePath(path string) string {
	return r.path + path
}
