package tests

import (
	"os"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg"
	fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	"github.com/onsi/gomega"

	"github.com/onsi/ginkgo/v2"
)

func ReportTest(report ginkgo.Report) {
	fernReporterBaseUrl := "http://localhost:8080/"

	// If FERN_REPORTER_BASE_URL is set, use it
	if os.Getenv("FERN_REPORTER_BASE_URL") != "" {
		fernReporterBaseUrl = os.Getenv("FERN_REPORTER_BASE_URL")
	}

	if os.Getenv("GITHUB_ACTION") == "" { //skip reporting in GH workflow
		fernApiClient := fern.New(pkg.PROJECT_ID, fern.WithBaseURL(fernReporterBaseUrl))
		err := fernApiClient.Report(report)
		gomega.Expect(err).To(gomega.BeNil(), "Unable to push report to Fern %v", err)
	}
}
