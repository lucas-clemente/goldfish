package notes_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNotes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Notes Suite")
}
