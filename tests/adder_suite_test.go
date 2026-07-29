package tests_test

import (
	"testing"

	"github.com/guidewire-oss/fern-ginkgo-client/v2/tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adder Suite", Label("this-is-a-suite-level-label"))
}

var _ = ReportAfterSuite("Fern reporting (direct client)", func(report Report) {
	tests.ReportTest(report)
})

var _ = ReportAfterSuite("Fern reporting (safe client)", func(report Report) {
	tests.ReportTestSafe(report)
})
