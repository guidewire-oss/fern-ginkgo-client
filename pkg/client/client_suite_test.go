package client_test

import (
	"github.com/guidewire-oss/fern-ginkgo-client/pkg"
	fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = ReportAfterSuite("", func(report Report) {
	if os.Getenv("GITHUB_ACTION") == "" { //skip reporting in GH workflow
		fern.ReportTestResult(pkg.ProjectName, report)
	}
})
