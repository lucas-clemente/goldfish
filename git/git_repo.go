package git

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	git2go "github.com/lucas-clemente/git2go"
	"github.com/lucas-clemente/treewatch"
)

// GitRepo is a git repository implementing the Repo interface for goldfish.
type GitRepo struct {
	path string
	repo *git2go.Repository
	tw   treewatch.TreeWatcher
	fo   *fanout
}

// NewGitRepo opens or makes a git repo at the given path
func NewGitRepo(repoPath string) (*GitRepo, error) {
	repoPath, err := filepath.EvalSymlinks(repoPath)
	if err != nil {
		return nil, err
	}

	repo, err := git2go.OpenRepository(repoPath)
	if err != nil {
		repo, err = git2go.InitRepository(repoPath, false)
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
		sig := &git2go.Signature{Name: "system", Email: "goldfish@clemente.io", When: time.Now()}
		_, err = repo.CreateCommit("refs/heads/master", sig, sig, "initial commit", tree)
		if err != nil {
			return nil, err
		}
	}

	tw, err := treewatch.NewTreeWatcher(repoPath)
	if err != nil {
		return nil, err
	}

	foChan := make(chan string)
	r := &GitRepo{path: repoPath, repo: repo, tw: tw, fo: newFanout(foChan)}

	go func() {
		for file := range tw.Changes() {
			if !strings.HasPrefix(file, r.path) {
				continue
			}
			file = strings.TrimPrefix(file, r.path)
			if strings.HasPrefix(file, "/.git") {
				continue
			}
			// Don't block commits on network
			go func() {
				foChan <- file
			}()
			log.Printf("file %s changed\n", file)
			err := r.addAllAndCommit("changed " + file)
			if err != nil {
				log.Println(err)
			}
		}
		close(foChan)
	}()

	return r, nil
}

// StopWatching stops watching for changes in the repo
func (r *GitRepo) StopWatching() {
	r.tw.Stop()
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

	_, err = io.Copy(file, data)
	return err
}

// ListFiles lists the files in a given directory
func (r *GitRepo) ListFiles(prefix string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(r.absolutePath(prefix))
	if err != nil {
		return nil, err
	}

	files := make([]string, 0, len(fileInfos))

	for _, f := range fileInfos {
		name := f.Name()
		if name[0] == '.' {
			continue
		}
		name = prefix + name
		if f.IsDir() {
			name += "/"
		}
		files = append(files, name)
	}

	return files, nil
}

// Observer sends file paths on changes
func (r *GitRepo) Observer() <-chan string {
	return r.fo.Output()
}

// CloseObserver closes an observer obtained from Observer()
func (r *GitRepo) CloseObserver(c <-chan string) {
	r.fo.Close(c)
}

func (r *GitRepo) addAllAndCommit(message string) error {
	index, err := r.repo.Index()
	if err != nil {
		return err
	}
	defer index.Free()

	if err := index.AddAll([]string{}, git2go.IndexAddDefault, nil); err != nil {
		return err
	}

	if err := index.UpdateAll([]string{}, nil); err != nil {
		return err
	}

	if err := index.Write(); err != nil {
		return err
	}

	treeID, err := index.WriteTree()
	if err != nil {
		return err
	}

	return r.commit(treeID, message)
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

	if *treeID == *headCommit.TreeId() {
		return nil
	}

	sig := &git2go.Signature{Name: "system", Email: "goldfish@clemente.io", When: time.Now()}
	_, err = r.repo.CreateCommit("refs/heads/master", sig, sig, message, tree, headCommit)
	if err != nil {
		return err
	}

	return nil
}

func (r *GitRepo) absolutePath(path string) string {
	return r.path + path
}
