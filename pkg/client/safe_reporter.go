package client

import (
	"fmt"
	"time"

	gt "github.com/onsi/ginkgo/v2/types"
)

// safeReportTimeout bounds how long ReportAfterSuiteSafe waits for Fern
// reporting to finish. It applies on top of (not instead of) whatever
// http.Client the caller supplies via ClientOption, so a caller-provided
// client with a zero Timeout can't block the suite forever. Declared as a
// var rather than a const so tests can shrink it to keep the timeout path
// fast to exercise.
var safeReportTimeout = 45 * time.Second

// ReportAfterSuiteSafe reports a Ginkgo suite result to Fern, logging a
// warning and returning normally on any failure, panic, or timeout instead
// of propagating it. Fern reporting is a side effect of running tests and
// must never fail or hang the suite it's reporting on.
func ReportAfterSuiteSafe(projectID string, report gt.Report, opts ...ClientOption) {
	done := make(chan struct{})

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("⚠️  Fern reporting failed: recovered from panic: %v\n", r)
			}
		}()

		c, err := New(projectID, opts...)
		if err != nil {
			fmt.Printf("⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
			return
		}

		if err := c.Report(report); err != nil {
			fmt.Printf("⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
		}
	}()

	select {
	case <-done:
	case <-time.After(safeReportTimeout):
		fmt.Printf("⚠️  Fern reporting failed: timed out after %s\n", safeReportTimeout)
	}
}
