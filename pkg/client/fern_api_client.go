package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// defaultHTTPTimeout bounds how long a Report() call can block on a stalled
// connection to Fern. Without it, a network stall never returns an error and
// the reporting suite hangs indefinitely instead of failing fast.
const defaultHTTPTimeout = 30 * time.Second

type FernApiClient struct {
	id           string
	httpClient   *http.Client
	baseURL      string
	token        string
	clientID     string
	clientSecret string
	authURL      string
	scope        string
}

type ClientOption func(*FernApiClient)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func New(projectId string, options ...ClientOption) (*FernApiClient, error) {
	return newWithContext(context.Background(), projectId, options...)
}

// newWithContext is the context-aware core of New. It exists so callers that
// need to bound the whole client-creation-plus-token-fetch operation (see
// ReportAfterSuiteSafe) can cancel the in-flight HTTP request on a deadline,
// rather than merely abandoning the wait for it.
func newWithContext(ctx context.Context, projectId string, options ...ClientOption) (*FernApiClient, error) {
	f := &FernApiClient{
		id:           projectId,
		httpClient:   &http.Client{Timeout: defaultHTTPTimeout},
		clientID:     os.Getenv("FERN_AUTH_CLIENT_ID"),
		clientSecret: os.Getenv("FERN_AUTH_CLIENT_SECRET"),
		authURL:      os.Getenv("AUTH_URL"),
		scope:        getEnvOrDefault("FERN_GINKGO_CLIENT_SCOPE", "fern.testrun.write"),
	}

	for _, o := range options {
		o(f)
	}

	// Generate token if credentials are available
	if f.clientID != "" && f.clientSecret != "" && f.authURL != "" {
		if err := f.generateToken(ctx); err != nil {
			return f, fmt.Errorf("failed to generate token: %w", err)
		}
	}

	return f, nil
}

func (f *FernApiClient) generateToken(ctx context.Context) error {
	tokenURL := strings.TrimRight(f.authURL, "/") + "/token"

	// Prepare the request body for client credentials flow
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", f.clientID)
	data.Set("client_secret", f.clientSecret)
	data.Set("scope", f.scope)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("token request failed with status: %d", resp.StatusCode)
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("failed to decode token response: %w", err)
	}

	f.token = tokenResp.AccessToken
	return nil
}

// GetToken returns the current access token
func (f *FernApiClient) GetToken() string {
	return f.token
}

// RefreshToken regenerates the access token
func (f *FernApiClient) RefreshToken() error {
	if f.clientID == "" || f.clientSecret == "" || f.authURL == "" {
		return fmt.Errorf("missing credentials for token refresh")
	}
	return f.generateToken(context.Background())
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

func WithCredentials(clientID, clientSecret, authURL string) ClientOption {
	return func(ac *FernApiClient) {
		ac.clientID = clientID
		ac.clientSecret = clientSecret
		ac.authURL = authURL
	}
}

func WithToken(token string) ClientOption {
	return func(ac *FernApiClient) {
		ac.token = token
	}
}

// GetBaseURL returns the configured base URL
func (f *FernApiClient) GetBaseURL() string {
	return f.baseURL
}

// GetClientID returns the configured client ID
func (f *FernApiClient) GetClientID() string {
	return f.clientID
}

// GetClientSecret returns the configured client secret
func (f *FernApiClient) getClientSecret() string {
	return f.clientSecret
}

// GetAuthURL returns the configured auth URL
func (f *FernApiClient) GetAuthURL() string {
	return f.authURL
}

// GetHTTPClient returns the configured http.Client
func (f *FernApiClient) GetHTTPClient() *http.Client {
	return f.httpClient
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
