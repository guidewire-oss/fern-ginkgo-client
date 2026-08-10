package client

import (
	"context"
	"fmt"
	"sync"
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

	// The worker goroutine and the select below can both decide a timeout
	// occurred (the worker sees ctx.Err() != nil after its request is
	// canceled; the select sees ctx.Done() fire) at nearly the same instant,
	// and `select` breaks that tie arbitrarily when both channels are ready.
	// warnOnce makes the outcome deterministic either way: exactly one
	// warning is emitted, never zero and never two.
	var warnOnce sync.Once
	warn := func(format string, args ...any) {
		warnOnce.Do(func() {
			fmt.Fprintf(ginkgo.GinkgoWriter, format, args...)
		})
	}
	timeoutMsg := fmt.Sprintf("⚠️  Fern reporting failed: timed out after %s\n", safeReportTimeout)

	go func() {
		defer close(done)
		defer func() {
			if r := recover(); r != nil {
				warn("⚠️  Fern reporting failed: recovered from panic: %v\n", r)
			}
		}()

		c, err := newWithContext(ctx, projectID, opts...)
		if err != nil {
			if ctx.Err() != nil {
				// The deadline had already elapsed by the time this
				// concluded; report it as a timeout, not a generic
				// creation failure.
				warn("%s", timeoutMsg)
			} else {
				warn("⚠️  Fern reporting failed: unable to create Fern API client: %v\n", err)
			}
			return
		}

		if err := c.reportWithContext(ctx, report); err != nil {
			if ctx.Err() != nil {
				warn("%s", timeoutMsg)
			} else {
				warn("⚠️  Fern reporting failed: unable to push report to Fern: %v\n", err)
			}
		}
	}()

	select {
	case <-done:
	case <-ctx.Done():
		warn("%s", timeoutMsg)
	}
}
