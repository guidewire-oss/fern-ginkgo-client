package client

import (
	"net/http"
	"time"
)

type FernApiClient struct {
	name          string
	httpClient    *http.Client
	baseURL       string
	username      string
	branch        string
	gitSHA        string
	project       string
	componentName string
}

type ClientOption func(*FernApiClient)

func New(testName string, options ...ClientOption) *FernApiClient {
	f := &FernApiClient{
		name:       testName,
		httpClient: http.DefaultClient,
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

func WithUsername(username string) ClientOption {
	return func(ac *FernApiClient) {
		ac.username = username
	}
}

func WithBranch(branch string) ClientOption {
	return func(ac *FernApiClient) {
		ac.branch = branch
	}
}

func WithGitSHA(sha string) ClientOption {
	return func(ac *FernApiClient) {
		ac.gitSHA = sha
	}
}

func WithProject(project string) ClientOption {
	return func(ac *FernApiClient) {
		ac.project = project
	}
}

func WithComponentName(component string) ClientOption {
	return func(ac *FernApiClient) {
		ac.componentName = component
	}
}
