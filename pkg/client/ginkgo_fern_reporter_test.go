package client

import (
	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("GinkGo Fern Reporter", Ordered, Label("unit"), func() {
	var testrun models.TestRun
	var originalEnv map[string]string

	BeforeEach(func() {
		testrun = models.TestRun{}

		// Save the original values to restore later
		originalEnv = map[string]string{
			"GITHUB_ACTION":           os.Getenv("GITHUB_ACTION"),
			"GIT_REPO_PATH":           os.Getenv("GIT_REPO_PATH"),
			"GITHUB_REF_NAME":         os.Getenv("GITHUB_REF_NAME"),
			"GITHUB_SHA":              os.Getenv("GITHUB_SHA"),
			"GITHUB_TRIGGERING_ACTOR": os.Getenv("GITHUB_TRIGGERING_ACTOR"),
			"GITHUB_SERVER_URL":       os.Getenv("GITHUB_SERVER_URL"),
			"GITHUB_REPOSITORY":       os.Getenv("GITHUB_REPOSITORY"),
			"GITHUB_RUN_ID":           os.Getenv("GITHUB_RUN_ID"),
		}
	})

	AfterEach(func() {
		// Restore original environment variables
		for key, val := range originalEnv {
			if val == "" {
				_ = os.Unsetenv(key) // Remove if it was originally unset
			} else {
				_ = os.Setenv(key, val) // Restore original value
			}
		}
	})

	It("should get local git details", func() {
		_ = os.Setenv("GITHUB_ACTION", "")
		addMetadataInfo(&testrun)

		Expect(testrun.GitBranch).ToNot(BeNil())
		Expect(len(testrun.GitBranch)).To(BeNumerically(">", 0))
	})

	It("should get local git details as NA", func() {
		_ = os.Setenv("GITHUB_ACTION", "")
		_ = os.Setenv("GIT_REPO_PATH", "/tmp/")
		addMetadataInfo(&testrun)

		Expect(testrun.GitBranch).To(Equal("NA"))
		Expect(testrun.GitSha).To(Equal("NA"))
	})

	It("should get GitHub git details", func() {

		_ = os.Setenv("GITHUB_ACTION", "true")
		_ = os.Setenv("GITHUB_REF_NAME", "feat/test")
		_ = os.Setenv("GITHUB_SHA", "acfb8356f058b88cef60a0b035df375f6471d6a0")
		_ = os.Setenv("GITHUB_TRIGGERING_ACTOR", "user")
		_ = os.Setenv("GITHUB_SERVER_URL", "https://github.com")
		_ = os.Setenv("GITHUB_REPOSITORY", "guidewire-oss/repo")
		_ = os.Setenv("GITHUB_RUN_ID", "13381986677")

		addMetadataInfo(&testrun)

		Expect(testrun.GitBranch).ToNot(BeNil())
		Expect(len(testrun.GitBranch)).To(BeNumerically(">", 0))
		Expect(testrun.GitBranch).To(Equal("feat/test"))
		Expect(testrun.GitSha).To(Equal("acfb8356f058b88cef60a0b035df375f6471d6a0"))
		Expect(testrun.BuildTriggerActor).To(Equal("user"))
		Expect(testrun.BuildUrl).To(Equal("https://github.com/guidewire-oss/repo/actions/runs/13381986677"))
	})
})
