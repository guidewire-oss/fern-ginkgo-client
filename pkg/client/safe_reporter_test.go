package client

import (
	"errors"
	"net/http"

	gt "github.com/onsi/ginkgo/v2/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReportAfterSuiteSafe", func() {
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
})
