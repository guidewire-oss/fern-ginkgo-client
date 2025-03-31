package tests_test

import (
	"github.com/guidewire-oss/fern-ginkgo-client/tests"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAdder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Adder Suite", Label("this-is-a-suite-level-label"))
}

var _ = ReportAfterSuite("", func(report Report) {
	tests.ReportTest(report)
})
