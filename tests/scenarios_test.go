package tests_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Reporting scenarios exist to dogfood the full range of SpecRun statuses
// Fern receives (passed, failed, skipped, pending) — the Adder specs above
// only ever produce "passed", so a real report-local run never exercised
// the others. Labeled "dogfood" (not "unit") so `make unit-test` — and CI,
// which only runs the "unit" label — never sees the intentional failure.
var _ = Describe("Reporting scenarios", Label("dogfood"), func() {
	It("passes", Label("scenario:pass"), func() {
		Expect(1 + 1).To(Equal(2))
	})

	It("fails", Label("scenario:fail"), func() {
		Expect(1 + 1).To(Equal(3))
	})

	It("is skipped", Label("scenario:skip"), func() {
		Skip("demonstrating a skipped SpecRun for Fern reporting")
	})

	PIt("is pending", Label("scenario:pending"), func() {
		Expect(true).To(BeFalse())
	})
})
