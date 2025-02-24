package client

import (
	"github.com/guidewire-oss/fern-ginkgo-client/pkg/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("GinkGo Fern Reporter", Ordered, Label("unit"), func() {
	var testrun models.TestRun

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

	BeforeEach(func() {
		testrun = models.TestRun{}
	})
})
