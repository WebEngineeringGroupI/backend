package grpc_test

import (
	"context"
	"errors"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gogrpc "google.golang.org/grpc"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Server", func() {
	var (
		connection      gogrpc.ClientConnInterface
		closeConnection context.CancelFunc
		urlValidator    *FakeURLValidator
	)
	BeforeEach(func() {
		urlValidator = &FakeURLValidator{returnValidURL: true}
		connection, closeConnection = newTestingConnection(grpc.Config{
			BaseDomain:         "https://example.com",
			ShortURLRepository: inmemory.NewRepository(),
			URLValidator:       urlValidator,
		})
	})

	AfterEach(func() {
		closeConnection()
	})

	Context("URLShorteningClient", func() {
		var client genproto.URLShorteningClient
		BeforeEach(func() {
			client = genproto.NewURLShorteningClient(connection)
		})

		It("returns the URLs shorted for an API call", func() {
			response, err := client.ShortURLs(context.Background(), &genproto.ShortURLsRequest{Urls: []string{"https://google.com", "https://youtube.com"}})

			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
			Expect(response.Results).To(HaveLen(2))
			Expect(response.Results[0].LongUrl).To(Equal("https://google.com"))
			Expect(response.Results[0].ShortUrl).To(Equal("https://example.com/r/cv6VxVdu"))
			Expect(response.Results[1].LongUrl).To(Equal("https://youtube.com"))
			Expect(response.Results[1].ShortUrl).To(Equal("https://example.com/r/unW6a4Dd"))
		})

		When("the URL is not valid", func() {
			It("returns the error", func() {
				urlValidator.shouldReturnValidURL(false)

				response, err := client.ShortURLs(context.Background(), &genproto.ShortURLsRequest{Urls: []string{"https://google.com", "https://youtube.com"}})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid long URL specified"))
				Expect(response).To(BeNil())
			})
		})

		When("the validator is unable to verify the URL", func() {
			It("returns the error", func() {
				urlValidator.shouldReturnError(errors.New("unknown testing error"))

				response, err := client.ShortURLs(context.Background(), &genproto.ShortURLsRequest{Urls: []string{"https://google.com", "https://youtube.com"}})

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown testing error"))
				Expect(response).To(BeNil())
			})
		})
	})
})
