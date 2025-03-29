package client

import (
	"net/http"
	"time"
)

type FernApiClient struct {
	name              string
	httpClient        *http.Client
	baseURL           string
	enableGeminiInsights bool
}

type ClientOption func(*FernApiClient)

func New(projectName string, options ...ClientOption) *FernApiClient {
	f := &FernApiClient{
		name:                 testName,
		httpClient:           http.DefaultClient,
		enableGeminiInsights: false, // Default to false
	}

	for _, o := range options {
		o(f)
	}

	return f
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(ac *FernApiClient) {
		ac.httpClient = httpClient
	}
}

func WithBaseURL(baseURL string) ClientOption {
	return func(ac *FernApiClient) {
		ac.baseURL = baseURL
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(ac *FernApiClient) {
		ac.httpClient.Timeout = timeout
	}
}

func WithGeminiInsights(enable bool) ClientOption {
	return func(ac *FernApiClient) {
		ac.enableGeminiInsights = enable
	}
}
