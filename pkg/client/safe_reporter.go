package client

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	gt "github.com/onsi/ginkgo/v2/types"
)

// safeReportTimeout bounds how long ReportAfterSuiteSafe waits for Fern
// reporting to finish. It applies on top of (not instead of) whatever
// http.Client the caller supplies via ClientOption: the deadline is carried
// on a context.Context threaded through to every HTTP request, so it
// actively cancels an in-flight request (rather than merely abandoning the
// wait for it) even when the caller's client has a zero Timeout. Declared
// as a var rather than a const so tests can shrink it to keep the timeout
// path fast to exercise.
var safeReportTimeout = 45 * time.Second

// ReportAfterSuiteSafe reports a Ginkgo suite result to Fern, logging a
// warning and returning normally on any failure, panic, or timeout instead
// of propagating it. Fern reporting is a side effect of running tests and
// must never fail or hang the suite it's reporting on.
func ReportAfterSuiteSafe(projectID string, report gt.Report, opts ...ClientOption) {
	ctx, cancel := context.WithTimeout(context.Background(), safeReportTimeout)
	defer cancel()

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(ginkgo.GinkgoWriter, "⚠️  Fern reporting failed: recovered from panic: %v\n", r)
			}
		}()

		c, err := newWithContext(ctx, projectID, opts...)
		if err != nil {
			// Once the deadline has fired, the caller already saw the
			// "timed out" message below; don't print a second, confusing
			// message for the same underlying cancellation.
			if ctx.Err() == nil {
				fmt.Fprintf(ginkgo.GinkgoWriter, "⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
			}
			return
		}

		if err := c.reportWithContext(ctx, report); err != nil {
			if ctx.Err() == nil {
				fmt.Fprintf(ginkgo.GinkgoWriter, "⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
			}
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		fmt.Fprintf(ginkgo.GinkgoWriter, "⚠️  Fern reporting failed: timed out after %s\n", safeReportTimeout)
	}
}
