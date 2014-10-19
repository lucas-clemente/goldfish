package notes_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/lucas-clemente/notes/notes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	Context("git repos", func() {
		var (
			tempDir string
			repo    notes.Repo
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "io.clemente.notes.test")
			Expect(err).To(BeNil())
			repo, err = notes.NewGitRepo(tempDir)
			Expect(err).To(BeNil())
			Expect(repo).ToNot(BeNil())
		})

		AfterEach(func() {
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

			cmd := exec.Command("git", "log")
			cmd.Dir = tempDir
			out, err := cmd.Output()
			Expect(err).To(BeNil())
			Expect(string(out)).To(ContainSubstring("Home.md"))
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
			Expect(err).To(MatchError(notes.NotFoundError{}))
		})
	})
})
