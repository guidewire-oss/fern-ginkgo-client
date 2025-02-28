package utils_test

import (
	"fmt"
	"github.com/guidewire-oss/fern-ginkgo-client/pkg/utils"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils Test", Label("unit"), func() {
	Describe("FindGitRoot", func() {
		var tempDir string

		BeforeEach(func() {
			var err error
			tempDir, err = os.MkdirTemp("", "test-git-root-*")
			Expect(err).ToNot(HaveOccurred())
		})

		AfterEach(func() {
			// Clean up the temporary directory
			err := os.RemoveAll(tempDir)
			Expect(err).ToNot(HaveOccurred(), "Failed to remove temp dir: %s", err)
		})

		It("should correctly find the git root when in a git subfolder", func() {
			// Create a fake .git directory at the repo root
			repoRoot := filepath.Join(tempDir, "repo")
			Expect(os.MkdirAll(filepath.Join(repoRoot, ".git"), 0755)).To(Succeed())

			// Create a subdirectory inside the repo
			subDir := filepath.Join(repoRoot, "subdir")
			Expect(os.Mkdir(subDir, 0755)).To(Succeed())

			// Call FindGitRoot from the subdirectory
			foundRoot, err := utils.FindGitRoot(subDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRoot).To(Equal(repoRoot))

			// Create another subdirectory inside the subdirectory
			subDir2 := filepath.Join(subDir, "subdir2")
			Expect(os.Mkdir(subDir2, 0755)).To(Succeed())

			// Call FindGitRoot from the subdirectory
			foundRoot, err = utils.FindGitRoot(subDir2)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRoot).To(Equal(repoRoot))

		})

		It("should return an error when no .git directory exists", func() {
			// Call FindGitRoot from a directory with no .git
			_, err := utils.FindGitRoot(tempDir)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("no valid .git directory found"))
		})

		It("should correctly detect a .git directory in the same folder", func() {
			// Create a .git directory in tempDir
			gitDir := filepath.Join(tempDir, ".git")
			Expect(os.Mkdir(gitDir, 0755)).To(Succeed())

			// Call FindGitRoot
			foundRoot, err := utils.FindGitRoot(tempDir)
			Expect(err).ToNot(HaveOccurred())
			Expect(foundRoot).To(Equal(tempDir))
		})

		It("should correctly find the git root when running in this folder with a relative path", func() {
			// Call FindGitRoot from this folder.
			currentFolder, err1 := utils.ToAbsolutePath(".")
			Expect(err1).ToNot(HaveOccurred())
			foundRoot, err2 := utils.FindGitRoot(".")
			Expect(err2).ToNot(HaveOccurred())
			Expect(foundRoot).ToNot(Equal(""))
			Expect(strings.HasPrefix(currentFolder, foundRoot)).To(Equal(true), fmt.Sprintf("Expected git root folder %s to be parent of current folder %s", foundRoot, currentFolder))
		})
	})
})
