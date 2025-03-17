package client_test

import (
	"github.com/guidewire-oss/fern-ginkgo-client/tests"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}

var _ = ReportAfterSuite("", func(report Report) {
	tests.ReportTest(report)
})
