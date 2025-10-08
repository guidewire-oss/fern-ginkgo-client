package client

import (
	"os"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ginkgo Fern Reporter", Ordered, Label("unit"), func() {
	var testrun models.TestRun
	var originalEnv map[string]string

	BeforeEach(func() {
		testrun = models.TestRun{}

		// Save original env values so we can restore them
		originalEnv = map[string]string{
			"GITHUB_ACTION":           os.Getenv("GITHUB_ACTION"),
			"GIT_REPO_PATH":           os.Getenv("GIT_REPO_PATH"),
			"GITHUB_REF_NAME":         os.Getenv("GITHUB_REF_NAME"),
			"GITHUB_SHA":              os.Getenv("GITHUB_SHA"),
			"GITHUB_TRIGGERING_ACTOR": os.Getenv("GITHUB_TRIGGERING_ACTOR"),
			"GITHUB_SERVER_URL":       os.Getenv("GITHUB_SERVER_URL"),
			"GITHUB_REPOSITORY":       os.Getenv("GITHUB_REPOSITORY"),
			"GITHUB_RUN_ID":           os.Getenv("GITHUB_RUN_ID"),
			"TEST_RUN_TAGS":           os.Getenv("TEST_RUN_TAGS"),
			"TEST_ENVIRONMENT":        os.Getenv("TEST_ENVIRONMENT"),
		}
	})

	AfterEach(func() {
		// Restore environment
		for key, val := range originalEnv {
			if val == "" {
				_ = os.Unsetenv(key)
			} else {
				_ = os.Setenv(key, val)
			}
		}
	})

	It("should get local git details", func() {
		_ = os.Setenv("GITHUB_ACTION", "")
		addMetadataInfo(&testrun)

		Expect(testrun.GitBranch).ToNot(BeEmpty())
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

		Expect(testrun.GitBranch).To(Equal("feat/test"))
		Expect(testrun.GitSha).To(Equal("acfb8356f058b88cef60a0b035df375f6471d6a0"))
		Expect(testrun.BuildTriggerActor).To(Equal("user"))
		Expect(testrun.BuildUrl).To(Equal("https://github.com/guidewire-oss/repo/actions/runs/13381986677"))
	})

	It("should set test run tags and environment from environment variables", func() {
		_ = os.Setenv("TEST_RUN_TAGS", "component:vector,team:qa")
		_ = os.Setenv("TEST_ENVIRONMENT", "staging")

		runTags := getRunLevelTags()
		Expect(runTags).To(HaveLen(2))
		Expect(runTags[0].Name).To(Equal("component:vector"))
		Expect(runTags[1].Name).To(Equal("team:qa"))

		testRun := models.TestRun{
			TestProjectID: "proj-1",
			TestSeed:      12345,
			Tags:          runTags,
			Environment:   os.Getenv("TEST_ENVIRONMENT"),
		}

		Expect(testRun.Tags).To(HaveLen(2))
		Expect(testRun.Environment).To(Equal("staging"))
	})
})
