package git

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"

	"github.com/lucas-clemente/goldfish/repository"
)

type gitFile struct {
	path    string
	modTime time.Time
	repo    *GitRepo
}

var _ repository.File = &gitFile{}

func (f *gitFile) Path() string {
	return f.path
}

func (f *gitFile) Reader() (io.ReadCloser, error) {
	return os.Open(f.repo.absolutePath(f.path))
}

func (f *gitFile) ModTime() time.Time {
	return f.modTime
}

// GitRepo is a git repository implementing the Repo interface for goldfish.
type GitRepo struct {
	path            string
	fo              *fanout
	notifyEventChan chan notify.EventInfo
}

var _ repository.Repo = &GitRepo{}

// NewGitRepo opens or makes a git repo at the given path
func NewGitRepo(repoPath string) (*GitRepo, error) {
	if _, err := os.Stat(repoPath + "/.git"); err != nil {
		// Run git init <dir>
		// Note that this creates the dir
		if err := exec.Command("git", "init", repoPath).Run(); err != nil {
			return nil, err
		}

		// Make an empty initial commit
		cmd := exec.Command("git", "commit", "--allow-empty", "-m", "initial commit")
		cmd.Dir = repoPath
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	}

	foChan := make(chan string)
	notifyEventChan := make(chan notify.EventInfo, 128)

	if err := notify.Watch(repoPath+"/...", notifyEventChan, notify.All); err != nil {
		return nil, err
	}

	r := &GitRepo{
		path:            repoPath,
		fo:              newFanout(foChan),
		notifyEventChan: notifyEventChan,
	}

	go func() {
		for eventInfo := range notifyEventChan {
			file := eventInfo.Path()

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
				log.Printf("error committing: %v\n", err)
			}
		}
		close(foChan)
	}()

	return r, nil
}

// LocalPathForFile returns the full local path for a file
func (r *GitRepo) LocalPathForFile(path string) (string, error) {
	return r.absolutePath(path), nil
}

// StopWatching stops watching for changes in the repo
func (r *GitRepo) StopWatching() {
	notify.Stop(r.notifyEventChan)
}

// ReadFile reads a file from the repo
func (r *GitRepo) ReadFile(path string) (repository.File, error) {
	info, err := os.Stat(r.absolutePath(path))
	if err != nil {
		return nil, err
	}
	return &gitFile{
		path:    path,
		modTime: info.ModTime(),
		repo:    r,
	}, nil
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

// DeleteFile deletes a file
func (r *GitRepo) DeleteFile(p string) error {
	return os.Remove(r.absolutePath(p))
}

// ListFiles lists the files in a given directory
func (r *GitRepo) ListFiles(prefix string) ([]repository.File, error) {
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}

	fileInfos, err := ioutil.ReadDir(r.absolutePath(prefix))
	if err != nil {
		return nil, err
	}

	files := make([]repository.File, 0, len(fileInfos))

	for _, f := range fileInfos {
		name := f.Name()
		if name[0] == '.' {
			continue
		}
		name = prefix + name
		if f.IsDir() {
			name += "/"
		}
		files = append(files, &gitFile{
			path:    name,
			modTime: f.ModTime(),
			repo:    r,
		})
	}

	return files, nil
}

// SearchFiles looks for markdown files containing `term` and returns the paths.
func (r *GitRepo) SearchFiles(term string) ([]repository.File, error) {
	term = strings.ToLower(term)

	// Walk through all files
	matches := []repository.File{}
	err := filepath.Walk(r.path, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "/.git/") || !strings.HasSuffix(path, ".md") {
			return nil
		}

		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		if strings.Contains(strings.ToLower(string(b)), term) {
			matches = append(matches, &gitFile{
				path:    strings.TrimPrefix(path, r.path),
				modTime: f.ModTime(),
				repo:    r,
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
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
	cmd := exec.Command("git", "add", "-A")
	cmd.Dir = r.path
	if err := cmd.Run(); err != nil {
		return err
	}

	// Check if there are changes
	cmd = exec.Command("git", "status", "--porcelain")
	cmd.Dir = r.path
	statusOutput, err := cmd.Output()
	if err != nil {
		return err
	}
	if len(statusOutput) == 0 {
		return nil
	}

	cmd = exec.Command("git", "commit", "-m", message)
	cmd.Dir = r.path
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (r *GitRepo) absolutePath(path string) string {
	return r.path + path
}
