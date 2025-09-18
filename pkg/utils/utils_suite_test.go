package utils_test

import (
	"testing"

	"github.com/guidewire-oss/fern-ginkgo-client/tests"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Utils Suite")
}

var _ = ReportAfterSuite("", func(report Report) {
	tests.ReportTest(report)
})
