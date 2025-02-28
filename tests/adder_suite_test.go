package tests_test

import (
	"github.com/guidewire-oss/fern-ginkgo-client/pkg"
	fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adder Suite", Label("this-is-a-suite-level-label"))
}

var _ = ReportAfterSuite("", func(report Report) {
	if os.Getenv("GITHUB_ACTION") == "" { //skip reporting in GH workflow
		fern.ReportTestResult(pkg.ProjectName, report)
	}
})
