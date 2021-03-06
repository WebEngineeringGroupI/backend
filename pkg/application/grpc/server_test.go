package grpc_test

import (
	"context"

	genproto "github.com/WebEngineeringGroupI/genproto-go/api/v1alpha1"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gogrpc "google.golang.org/grpc"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	urlmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/eventstore/inmemory"
)

var _ = Describe("Server", func() {
	var (
		ctrl                       *gomock.Controller
		metrics                    *urlmocks.MockMetrics
		connection                 gogrpc.ClientConnInterface
		closeConnection            context.CancelFunc
		shortURLRepository         event.Repository
		loadBalancerURLsRepository event.Repository
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		metrics = urlmocks.NewMockMetrics(ctrl)
		shortURLRepository = event.NewRepository(&url.ShortURL{}, inmemory.NewEventStore(), event.NewBroker())
		loadBalancerURLsRepository = event.NewRepository(&url.LoadBalancedURL{}, inmemory.NewEventStore(), event.NewBroker())

		connection, closeConnection = newTestingConnection(grpc.Config{
			BaseDomain:                 "https://example.com",
			CustomMetrics:              metrics,
			ShortURLRepository:         shortURLRepository,
			LoadBalancedURLsRepository: loadBalancerURLsRepository,
		})
		metrics.EXPECT().RecordSingleURLMetrics().AnyTimes()
		metrics.EXPECT().RecordFileURLMetrics().AnyTimes()
	})

	AfterEach(func() {
		closeConnection()
		ctrl.Finish()
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

		It("shortens a single URL", func() {
			response, err := client.ShortSingleURL(ctx, &genproto.ShortSingleURLRequest{
				Url: "https://google.com",
			})

			Expect(err).ToNot(HaveOccurred())
			Expect(response.GetShortUrl()).To(Equal("https://example.com/r/cv6VxVdu"))
			Expect(response.GetLongUrl()).To(Equal("https://google.com"))
		})
		When("the url is empty", func() {
			It("returns an error", func() {
				response, err := client.ShortSingleURL(ctx, &genproto.ShortSingleURLRequest{})

				Expect(err).To(MatchError(ContainSubstring("empty URL provided")))
				Expect(response).To(BeNil())
			})
		})
	})
})
