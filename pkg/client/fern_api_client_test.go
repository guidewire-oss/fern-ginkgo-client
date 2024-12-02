package client_test

import (
	"time"

	"github.com/guidewire-oss/fern-ginkgo-client/pkg/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("FernApiClient", func() {
	It("should get a new client", func() {

		fernApiClient := client.New("test")

		Expect(fernApiClient).ToNot(BeNil())

	})

	It("should get a new client with BaseURL", func() {

		fernApiClient := client.New("test", client.WithBaseURL("test URL"))

		Expect(fernApiClient).ToNot(BeNil())

	})

	It("should get a new client with HTTP Client", func() {

		fernApiClient := client.New("test", client.WithHTTPClient(nil))

		Expect(fernApiClient).ToNot(BeNil())

	})

	It("should get a new client with timeout", func() {

		fernApiClient := client.New("test", client.WithTimeout(5*time.Second))

		Expect(fernApiClient).ToNot(BeNil())

	})

	It("should get a new client with Gemini Insights enabled", func() {
		
		fernApiClient := client.New("test", client.WithGeminiInsights(true))

		Expect(fernApiClient).ToNot(BeNil())
		// We need to add a method to check if Gemini Insights is enabled
		// For now, we'll just check that the client is created
	})
})
