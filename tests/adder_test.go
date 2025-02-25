package tests_test

import (
	fern "github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"

	. "github.com/guidewire-oss/fern-ginkgo-client/tests"
)

var _ = Describe("Adder", func() {

	Describe("Add", func() {

		It("adds two numbers", func() {
			sum := Add(2, 3)
			Expect(sum).To(Equal(5))
		})
	})

})
var _ = ReportAfterSuite("", func(report Report) {
	if os.Getenv("GITHUB_ACTION") == "" { //skip reporting in GH workflow
		f := fern.New("Example Test",
			fern.WithBaseURL("http://localhost:8080/"),
		)

		err := f.Report("example test", report)

		Expect(err).To(BeNil(), "Unable to create reporter file")
	}
})
