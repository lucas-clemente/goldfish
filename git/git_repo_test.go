package git_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/lucas-clemente/goldfish/git"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExpectSoon(f func() bool) {
	for i := 0; i < 100; i++ {
		if f() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	Expect(f()).To(BeTrue())
}

var _ = Describe("Repo", func() {
	Context("git repos", func() {
		var (
			tempDir string
			repo    *git.GitRepo
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "io.clemente.git.test")
			Expect(err).To(BeNil())
			repo, err = git.NewGitRepo(tempDir)
			Expect(err).To(BeNil())
			Expect(repo).ToNot(BeNil())
		})

		AfterEach(func() {
			repo.StopWatching()
			// Give the fs events some time to get processed before deleting the repo
			time.Sleep(200 * time.Millisecond)
			os.RemoveAll(tempDir)
		})

		It("inits new repos", func() {
			cmd := exec.Command("git", "log")
			cmd.Dir = tempDir
			out, err := cmd.Output()
			Expect(err).To(BeNil())
			Expect(string(out)).To(ContainSubstring("initial commit"))
		})

		It("saves and reads files", func() {
			err := repo.StoreFile("/foo/Home.md", bytes.NewBufferString("foobar"))
			Expect(err).To(BeNil())
			reader, err := repo.ReadFile("/foo/Home.md")
			Expect(err).To(BeNil())
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("foobar")))

			ExpectSoon(func() bool {
				cmd := exec.Command("git", "log")
				cmd.Dir = tempDir
				out, err := cmd.Output()
				Expect(err).To(BeNil())
				return strings.Contains(string(out), "foo")
			})
		})

		It("updates and reads files", func() {
			err := repo.StoreFile("/foo/Home.md", bytes.NewBufferString("foobar"))
			Expect(err).To(BeNil())

			err = repo.StoreFile("/foo/Home.md", bytes.NewBufferString("foobaz"))
			Expect(err).To(BeNil())

			reader, err := repo.ReadFile("/foo/Home.md")
			Expect(err).To(BeNil())
			defer reader.Close()
			data, err := ioutil.ReadAll(reader)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte("foobaz")))

			ExpectSoon(func() bool {
				cmd := exec.Command("git", "log")
				cmd.Dir = tempDir
				out, err := cmd.Output()
				Expect(err).To(BeNil())
				return strings.Contains(string(out), "foo")
			})
		})

		It("lists files", func() {
			err := repo.StoreFile("/foo/bar", bytes.NewBufferString("foobar"))
			Expect(err).To(BeNil())

			err = repo.StoreFile("/baz", bytes.NewBufferString("foobaz"))
			Expect(err).To(BeNil())

			files, err := repo.ListFiles("/")
			Expect(err).To(BeNil())
			Expect(files).To(Equal([]string{"/baz", "/foo/"}))

			files, err = repo.ListFiles("/foo/")
			Expect(err).To(BeNil())
			Expect(files).To(Equal([]string{"/foo/bar"}))
		})

		It("handles not found", func() {
			reader, err := repo.ReadFile("/foo")
			Expect(reader).To(BeNil())
			Expect(os.IsNotExist(err)).To(BeTrue())
		})
	})
})
