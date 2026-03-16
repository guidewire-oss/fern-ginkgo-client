package tests

import (
	"fmt"
	"os"

	"github.com/guidewire-oss/fern-ginkgo-client/v2/pkg"
	fern "github.com/guidewire-oss/fern-ginkgo-client/v2/pkg/client"

	. "github.com/onsi/ginkgo/v2"
)

func ReportTest(report Report) {
	fernBaseUrl := "http://localhost:8080/"
	fernProjectId := pkg.PROJECT_ID

	// If FERN_BASE_URL is set, use it
	if os.Getenv("FERN_BASE_URL") != "" {
		fernBaseUrl = os.Getenv("FERN_BASE_URL")
	}

	if os.Getenv("PROJECT_ID") != "" {
		fernProjectId = os.Getenv("PROJECT_ID")
	}

	if os.Getenv("GITHUB_ACTION") == "" { //skip reporting in GH workflow
		fernApiClient, err := fern.New(fernProjectId, fern.WithBaseURL(fernBaseUrl))
		if err != nil {
			fmt.Printf("⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
			return
		}
		err = fernApiClient.Report(report)
		if err != nil {
			fmt.Printf("⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
		}
	}
}
