package notes_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	. "github.com/lucas-clemente/notes/notes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo", func() {
	Context("git repos", func() {
		var (
			tempDir string
		)

		BeforeEach(func() {
			var err error
			tempDir, err = ioutil.TempDir("", "io.clemente.notes.test")
			Expect(err).To(BeNil())
		})

		AfterEach(func() {
			os.RemoveAll(tempDir)
		})

		It("inits new repos", func() {
			repo, err := NewGitRepo(tempDir)
			Expect(err).To(BeNil())
			Expect(repo).ToNot(BeNil())

			cmd := exec.Command("git", "log")
			cmd.Dir = tempDir
			out, err := cmd.Output()
			Expect(err).To(BeNil())
			Expect(string(out)).To(ContainSubstring("initial commit"))
		})

		It("saves and reads files", func() {
			d := []byte("# Home\n")
			repo, err := NewGitRepo(tempDir)
			Expect(err).To(BeNil())
			Expect(repo).ToNot(BeNil())
			err = repo.StoreFile("/foo/Home.md", d)
			Expect(err).To(BeNil())
			data, err := repo.ReadFile("/foo/Home.md")
			Expect(err).To(BeNil())
			Expect(data).To(Equal(d))

			cmd := exec.Command("git", "log")
			cmd.Dir = tempDir
			out, err := cmd.Output()
			Expect(err).To(BeNil())
			Expect(string(out)).To(ContainSubstring("Home.md"))
		})
	})
})
