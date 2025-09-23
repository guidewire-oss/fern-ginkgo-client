package client

import (
	_ "encoding/json"
	"errors"
	"io"
	"net/http"
	_ "net/http/httptest"
	"os"
	_ "os"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// mockRoundTripper implements http.RoundTripper for mocking requests
type mockRoundTripper struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

var _ = Describe("FernApiClient", func() {
	var mockHTTPClient *http.Client

	var orig_client_id, orig_client_secret string

	BeforeEach(func() {
		mockHTTPClient = &http.Client{
			Transport: &mockRoundTripper{},
		}

		orig_client_id = os.Getenv("FERN_AUTH_CLIENT_ID")
		orig_client_secret = os.Getenv("FERN_AUTH_CLIENT_SECRET")
		_ = os.Unsetenv("FERN_AUTH_CLIENT_ID")
		_ = os.Unsetenv("FERN_AUTH_CLIENT_SECRET")
		DeferCleanup(func() {
			_ = os.Setenv("FERN_AUTH_CLIENT_ID", orig_client_id)
			_ = os.Setenv("FERN_AUTH_CLIENT_SECRET", orig_client_secret)
		})
	})

	Describe("New", func() {
		It("creates client without credentials", func() {
			c, err := New("proj1", WithHTTPClient(mockHTTPClient))
			Expect(err).To(BeNil())
			Expect(c).NotTo(BeNil())
			Expect(c.GetToken()).To(Equal(""))
		})

		It("returns error when token generation fails", func() {
			rt := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("network error")
				},
			}
			c, err := New("proj2",
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
			Expect(err).To(HaveOccurred())
			Expect(c).NotTo(BeNil())
		})
	})

	Describe("generateToken", func() {
		It("successfully sets token on valid response", func() {
			tokenJSON := `{"access_token":"abc123","token_type":"bearer","expires_in":3600}`
			rt := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					Expect(req.Method).To(Equal("POST"))
					Expect(req.URL.String()).To(Equal("http://auth/token"))
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(tokenJSON)),
					}, nil
				},
			}
			c, err := New("proj3",
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
			Expect(err).To(BeNil())
			Expect(c.GetToken()).To(Equal("abc123"))
		})

		It("returns error on http failure", func() {
			rt := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return nil, errors.New("boom")
				},
			}
			c, _ := New("projX",
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
			err := c.RefreshToken()
			Expect(err).To(HaveOccurred())
		})

		It("returns error on non-200 response", func() {
			rt := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 400,
						Body:       io.NopCloser(strings.NewReader("bad")),
					}, nil
				},
			}
			_, err := New("projY",
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
			Expect(err).To(HaveOccurred())
		})

		It("returns error on invalid JSON", func() {
			rt := &mockRoundTripper{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("not-json")),
					}, nil
				},
			}
			c, _ := New("projZ",
				WithHTTPClient(&http.Client{Transport: rt}),
				WithCredentials("id", "secret", "http://auth"),
			)
			err := c.RefreshToken()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("RefreshToken", func() {
		It("returns error when missing credentials", func() {
			c := FernApiClient{}
			err := c.RefreshToken()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Options", func() {
		It("sets baseURL", func() {
			c, _ := New("proj", WithBaseURL("http://api"))
			Expect(c.GetBaseURL()).To(Equal("http://api"))
		})

		It("sets timeout", func() {
			c, _ := New("proj", WithTimeout(2*time.Second))
			Expect(c.GetHTTPClient().Timeout).To(Equal(2 * time.Second))
		})

		It("sets token", func() {
			c, _ := New("proj", WithToken("tok123"))
			Expect(c.GetToken()).To(Equal("tok123"))
		})

		It("sets credentials", func() {
			c, _ := New("proj", WithCredentials("id", "secret", "http://auth"))
			Expect(c.GetClientID()).To(Equal("id"))
			Expect(c.getClientSecret()).To(Equal("secret"))
			Expect(c.GetAuthURL()).To(Equal("http://auth"))
		})
	})
})

var _ = Describe("getEnvOrDefault", func() {
	const (
		testEnvKey     = "TEST_ENV_VAR"
		testEnvValue   = "test_value"
		defaultValue   = "default_value"
		emptyValue     = ""
		anotherKey     = "ANOTHER_TEST_KEY"
		anotherValue   = "another_value"
		anotherDefault = "another_default"
	)

	// Store original environment state for cleanup
	var originalEnvValue string
	var originalEnvExists bool

	BeforeEach(func() {
		// Store the original state of our test environment variable
		originalEnvValue, originalEnvExists = os.LookupEnv(testEnvKey)

		// Clean up any existing test environment variables
		os.Unsetenv(testEnvKey)
		os.Unsetenv(anotherKey)
	})

	AfterEach(func() {
		// Restore original environment state
		if originalEnvExists {
			os.Setenv(testEnvKey, originalEnvValue)
		} else {
			os.Unsetenv(testEnvKey)
		}

		// Clean up test environment variables
		os.Unsetenv(anotherKey)
	})

	Describe("when environment variable exists and has a value", func() {
		BeforeEach(func() {
			os.Setenv(testEnvKey, testEnvValue)
		})

		It("should return the environment variable value", func() {
			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).To(Equal(testEnvValue))
		})

		It("should not return the default value", func() {
			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).NotTo(Equal(defaultValue))
		})
	})

	Describe("when environment variable does not exist", func() {
		It("should return the default value", func() {
			// Ensure the environment variable is not set
			os.Unsetenv(testEnvKey)

			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).To(Equal(defaultValue))
		})
	})

	Describe("when environment variable exists but is empty", func() {
		BeforeEach(func() {
			os.Setenv(testEnvKey, emptyValue)
		})

		It("should return the default value", func() {
			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).To(Equal(defaultValue))
		})

		It("should not return an empty string", func() {
			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("when default value is empty", func() {
		Context("and environment variable exists with a value", func() {
			BeforeEach(func() {
				os.Setenv(testEnvKey, testEnvValue)
			})

			It("should return the environment variable value", func() {
				result := getEnvOrDefault(testEnvKey, emptyValue)
				Expect(result).To(Equal(testEnvValue))
			})
		})

		Context("and environment variable does not exist", func() {
			It("should return the empty default value", func() {
				os.Unsetenv(testEnvKey)

				result := getEnvOrDefault(testEnvKey, emptyValue)
				Expect(result).To(Equal(emptyValue))
				Expect(result).To(BeEmpty())
			})
		})

		Context("and environment variable exists but is empty", func() {
			BeforeEach(func() {
				os.Setenv(testEnvKey, emptyValue)
			})

			It("should return the empty default value", func() {
				result := getEnvOrDefault(testEnvKey, emptyValue)
				Expect(result).To(Equal(emptyValue))
				Expect(result).To(BeEmpty())
			})
		})
	})

	Describe("edge cases", func() {
		Context("when environment variable contains only whitespace", func() {
			BeforeEach(func() {
				os.Setenv(testEnvKey, "   ")
			})

			It("should return the whitespace value, not the default", func() {
				result := getEnvOrDefault(testEnvKey, defaultValue)
				Expect(result).To(Equal("   "))
				Expect(result).NotTo(Equal(defaultValue))
			})
		})

		Context("when environment variable contains special characters", func() {
			specialValue := "special!@#$%^&*()_+-={}[]|\\:;\"'<>?,./"

			BeforeEach(func() {
				os.Setenv(testEnvKey, specialValue)
			})

			It("should return the special characters value", func() {
				result := getEnvOrDefault(testEnvKey, defaultValue)
				Expect(result).To(Equal(specialValue))
			})
		})

		Context("when environment variable contains newlines", func() {
			multilineValue := "line1\nline2\nline3"

			BeforeEach(func() {
				os.Setenv(testEnvKey, multilineValue)
			})

			It("should return the multiline value", func() {
				result := getEnvOrDefault(testEnvKey, defaultValue)
				Expect(result).To(Equal(multilineValue))
			})
		})
	})

	Describe("multiple calls", func() {
		It("should return consistent results for the same inputs", func() {
			os.Setenv(testEnvKey, testEnvValue)

			result1 := getEnvOrDefault(testEnvKey, defaultValue)
			result2 := getEnvOrDefault(testEnvKey, defaultValue)

			Expect(result1).To(Equal(result2))
			Expect(result1).To(Equal(testEnvValue))
		})

		It("should handle different environment variables independently", func() {
			os.Setenv(testEnvKey, testEnvValue)
			os.Setenv(anotherKey, anotherValue)

			result1 := getEnvOrDefault(testEnvKey, defaultValue)
			result2 := getEnvOrDefault(anotherKey, anotherDefault)

			Expect(result1).To(Equal(testEnvValue))
			Expect(result2).To(Equal(anotherValue))
		})
	})

	Describe("function signature and behavior", func() {
		It("should accept string parameters and return string", func() {
			result := getEnvOrDefault(testEnvKey, defaultValue)
			Expect(result).To(BeAssignableToTypeOf(""))
		})

		It("should not modify environment variables", func() {
			originalValue := "original"
			os.Setenv(testEnvKey, originalValue)

			getEnvOrDefault(testEnvKey, defaultValue)

			currentValue := os.Getenv(testEnvKey)
			Expect(currentValue).To(Equal(originalValue))
		})
	})
})
