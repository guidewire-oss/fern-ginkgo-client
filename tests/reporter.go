package tests

import (
	"fmt"
	"os"

	"github.com/guidewire-oss/fern-ginkgo-client/v2/pkg"
	fern "github.com/guidewire-oss/fern-ginkgo-client/v2/pkg/client"

	. "github.com/onsi/ginkgo/v2"
)

// resolveFernTarget reads the local Fern target from env vars, falling
// back to this project's own dev defaults.
func resolveFernTarget() (baseURL, projectID string) {
	baseURL = "http://localhost:8080/"
	if v := os.Getenv("FERN_BASE_URL"); v != "" {
		baseURL = v
	}

	projectID = pkg.PROJECT_ID
	if v := os.Getenv("PROJECT_ID"); v != "" {
		projectID = v
	}
	return baseURL, projectID
}

// ReportTest reports via the direct New()+Report() pattern most existing
// consumers (e.g. ccs-atmos-tests) still use today. The suite name is
// tagged "(direct client)" so this run is distinguishable from ReportTestSafe's
// in Fern.
func ReportTest(report Report) {
	if os.Getenv("GITHUB_ACTION") != "" { //skip reporting in GH workflow
		return
	}

	report.SuiteDescription += " (direct client)"

	baseURL, projectID := resolveFernTarget()
	fernApiClient, err := fern.New(projectID, fern.WithBaseURL(baseURL))
	if err != nil {
		fmt.Printf("⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
		return
	}
	if err := fernApiClient.Report(report); err != nil {
		fmt.Printf("⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
	}
}

// ReportTestSafe reports the same suite via ReportAfterSuiteSafe, the
// safe-by-default helper. The suite name is tagged "(safe client)" so this
// run is distinguishable from ReportTest's in Fern.
func ReportTestSafe(report Report) {
	if os.Getenv("GITHUB_ACTION") != "" { //skip reporting in GH workflow
		return
	}

	report.SuiteDescription += " (safe client)"

	baseURL, projectID := resolveFernTarget()
	fern.ReportAfterSuiteSafe(projectID, report, fern.WithBaseURL(baseURL))
}
