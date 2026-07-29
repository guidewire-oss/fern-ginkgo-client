package client

import (
	"fmt"

	gt "github.com/onsi/ginkgo/v2/types"
)

// ReportAfterSuiteSafe reports a Ginkgo suite result to Fern, logging a
// warning and returning normally on any failure instead of propagating the
// error. Fern reporting is a side effect of running tests and must never
// fail or hang the suite it's reporting on.
func ReportAfterSuiteSafe(projectID string, report gt.Report, opts ...ClientOption) {
	c, err := New(projectID, opts...)
	if err != nil {
		fmt.Printf("⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
		return
	}

	if err := c.Report(report); err != nil {
		fmt.Printf("⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
	}
}
