package grpc_test

import (
	"context"

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
		metrics         *FakeMetrics
	)
	BeforeEach(func() {
		metrics = &FakeMetrics{}
		connection, closeConnection = newTestingConnection(grpc.Config{
			BaseDomain:                 "https://example.com",
			CustomMetrics:              metrics,
			ShortURLRepository:         inmemory.NewRepository(),
			LoadBalancedURLsRepository: inmemory.NewRepository(),
			EventOutbox:                inmemory.NewRepository(),
		})
	})

	AfterEach(func() {
		closeConnection()
	})

	Context("URLShorteningClient", func() {
		var (
			ctx             context.Context
			client          genproto.URLShorteningClient
			shortURLsClient genproto.URLShortening_ShortURLsClient
		)
		BeforeEach(func() {
			ctx = context.Background()
			client = genproto.NewURLShorteningClient(connection)
			var err error
			shortURLsClient, err = client.ShortURLs(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns the URLs shorted for an API call", func() {
			err := shortURLsClient.Send(&genproto.ShortURLsRequest{Url: "https://google.com"})
			Expect(err).ToNot(HaveOccurred())
			err = shortURLsClient.Send(&genproto.ShortURLsRequest{Url: "https://youtube.com"})
			Expect(err).ToNot(HaveOccurred())

			response, err := shortURLsClient.Recv()
			Expect(err).ToNot(HaveOccurred())
			Expect(response).ToNot(BeNil())
			Expect(response.GetSuccess()).ToNot(BeNil())
			Expect(response.GetSuccess().LongUrl).To(Equal("https://google.com"))
			Expect(response.GetSuccess().ShortUrl).To(Equal("https://example.com/r/cv6VxVdu"))

			response, err = shortURLsClient.Recv()
			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetSuccess()).ToNot(BeNil())
			Expect(response.GetSuccess().LongUrl).To(Equal("https://youtube.com"))
			Expect(response.GetSuccess().ShortUrl).To(Equal("https://example.com/r/unW6a4Dd"))
		})

		When("the client wants to create a load-balanced URL", func() {
			It("creates the URL correctly", func() {
				balanceURLsResponse, err := client.BalanceURLs(ctx, &genproto.BalanceURLsRequest{Urls: []string{"https://google.com", "https://youtube.com"}})

				Expect(err).ToNot(HaveOccurred())
				Expect(balanceURLsResponse.GetShortUrl()).To(Equal("https://example.com/lb/8YOPuCnc"))
			})
			Context("but the list is empty", func() {
				It("returns an error", func() {
					balanceURLsResponse, err := client.BalanceURLs(ctx, &genproto.BalanceURLsRequest{})

					Expect(err).To(MatchError(ContainSubstring("no URLs specified")))
					Expect(balanceURLsResponse.GetShortUrl()).To(BeEmpty())
				})
			})
		})
	})
})
