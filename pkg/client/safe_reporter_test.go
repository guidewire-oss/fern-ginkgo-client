package client

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	gt "github.com/onsi/ginkgo/v2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReportAfterSuiteSafe", func() {
	var origClientID, origClientSecret, origAuthURL string

	BeforeEach(func() {
		// Credential env vars change which path New() takes (token
		// generation vs. skip). Pin them so each test exercises the path
		// its ClientOptions describe, regardless of the ambient environment.
		origClientID = os.Getenv("FERN_AUTH_CLIENT_ID")
		origClientSecret = os.Getenv("FERN_AUTH_CLIENT_SECRET")
		origAuthURL = os.Getenv("AUTH_URL")
		_ = os.Unsetenv("FERN_AUTH_CLIENT_ID")
		_ = os.Unsetenv("FERN_AUTH_CLIENT_SECRET")
		_ = os.Unsetenv("AUTH_URL")
		DeferCleanup(func() {
			_ = os.Setenv("FERN_AUTH_CLIENT_ID", origClientID)
			_ = os.Setenv("FERN_AUTH_CLIENT_SECRET", origClientSecret)
			_ = os.Setenv("AUTH_URL", origAuthURL)
		})
	})

	It("does not panic when client creation fails", func() {
		rt := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network error")
			},
		}

		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-1", gt.Report{},
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
		}).NotTo(Panic())
	})

	It("does not panic when pushing the report fails", func() {
		rt := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("connection refused")
			},
		}

		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-2", gt.Report{},
				WithHTTPClient(&http.Client{Transport: rt}),
				WithBaseURL("http://fern.invalid"),
			)
		}).NotTo(Panic())
	})

	It("does not panic when reporting succeeds", func() {
		rt := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 200,
					Body:       http.NoBody,
				}, nil
			},
		}

		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-3", gt.Report{},
				WithHTTPClient(&http.Client{Transport: rt}),
				WithBaseURL("http://fern.invalid"),
			)
		}).NotTo(Panic())
	})

	It("does not panic when the caller supplies a nil http.Client", func() {
		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-nil-client", gt.Report{},
				WithHTTPClient(nil),
				WithBaseURL("http://fern.invalid"),
			)
		}).NotTo(Panic())
	})

	It("cancels the in-flight request via context when the deadline fires, rather than only abandoning the wait", func() {
		origTimeout := safeReportTimeout
		safeReportTimeout = 50 * time.Millisecond
		DeferCleanup(func() {
			safeReportTimeout = origTimeout
		})

		// A well-behaved RoundTripper aborts as soon as its request's
		// context is done. Blocking on that (instead of on an
		// independent channel) proves ReportAfterSuiteSafe's deadline
		// actually cancels the request rather than merely giving up on
		// waiting for it.
		roundTripErr := make(chan error, 1)
		rt := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				<-req.Context().Done()
				err := req.Context().Err()
				roundTripErr <- err
				return nil, err
			},
		}

		start := time.Now()
		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-cancel", gt.Report{},
				WithHTTPClient(&http.Client{Transport: rt}), // Timeout left at zero (no timeout)
				WithBaseURL("http://fern.invalid"),
			)
		}).NotTo(Panic())
		elapsed := time.Since(start)

		// It must have genuinely waited out the deadline...
		Expect(elapsed).To(BeNumerically(">=", 50*time.Millisecond))
		// ...but not much longer than it, proving the wait was bounded.
		Expect(elapsed).To(BeNumerically("<", 2*time.Second))

		// And the request's context must have actually been canceled by
		// the deadline, not left running in the background.
		var rtErr error
		Eventually(roundTripErr, time.Second).Should(Receive(&rtErr))
		Expect(rtErr).To(MatchError(context.DeadlineExceeded))
	})

	It("does not hang when a stalled connection ignores context cancellation", func() {
		origTimeout := safeReportTimeout
		safeReportTimeout = 50 * time.Millisecond
		DeferCleanup(func() {
			safeReportTimeout = origTimeout
		})

		// Some RoundTrippers (or the transports they wrap) don't honor
		// context cancellation. ReportAfterSuiteSafe must still return
		// promptly on the deadline even though the background goroutine
		// is left to leak until blockUntil is eventually closed below.
		blockUntil := make(chan struct{})
		DeferCleanup(func() { close(blockUntil) })

		rt := &mockRoundTripper{
			roundTripFunc: func(req *http.Request) (*http.Response, error) {
				<-blockUntil // simulate a stalled connection that never responds
				return nil, errors.New("unreachable")
			},
		}

		start := time.Now()
		Expect(func() {
			ReportAfterSuiteSafe("proj-safe-stall", gt.Report{},
				WithHTTPClient(&http.Client{Transport: rt}), // Timeout left at zero (no timeout)
				WithBaseURL("http://fern.invalid"),
			)
		}).NotTo(Panic())
		elapsed := time.Since(start)

		// The stalled round-trip never returns on its own; the only way
		// this completes at all is the safeReportTimeout deadline firing,
		// and it must not fire early.
		Expect(elapsed).To(BeNumerically(">=", 50*time.Millisecond))
		Expect(elapsed).To(BeNumerically("<", 2*time.Second))
	})
})
