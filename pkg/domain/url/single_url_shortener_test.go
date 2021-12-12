package url_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	domainmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Single URL shortener", func() {
	var (
		controller *gomock.Controller
		shortener  *url.SingleURLShortener
		repository url.ShortURLRepository
		metrics    *mocks.MockMetrics
		emitter    *domainmocks.MockEmitter
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repository = inmemory.NewRepository()

		controller = gomock.NewController(GinkgoT())
		emitter = domainmocks.NewMockEmitter(controller)
		metrics = mocks.NewMockMetrics(controller)

		shortener = url.NewSingleURLShortener(repository, metrics, emitter)
	})

	AfterEach(func() {
		controller.Finish()
	})

	Context("when providing a long URL", func() {
		BeforeEach(func() {
			emitter.EXPECT().EmitShortURLCreated(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
		})

		It("generates a hash", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.Hash).To(Equal("cv6VxVdu"))
		})

		It("contains the real value from the original URL", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.OriginalURL.URL).To(Equal("https://google.com"))
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				metrics.EXPECT().RecordSingleURLMetrics().Times(2)
				shortGoogleURL, err := shortener.HashFromURL(ctx, "https://google.com")
				Expect(err).ToNot(HaveOccurred())

				shortFacebookURL, err := shortener.HashFromURL(ctx, "https://facebook.com")
				Expect(err).ToNot(HaveOccurred())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
			})
		})

		It("stores the short URL in a repository", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			shortURL, err := shortener.HashFromURL(ctx, "https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.Hash).To(Equal(shortURL.Hash))
		})

		It("stores the URL as non verified", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			shortURL, err := shortener.HashFromURL(ctx, "https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.OriginalURL.IsValid).To(BeFalse())
		})

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})

	It("should emit an event saved to the emitter", func() {
		metrics.EXPECT().RecordSingleURLMetrics().Times(1)
		emitter.EXPECT().EmitShortURLCreated(gomock.Any(), "cv6VxVdu", "https://google.com")
		_, err := shortener.HashFromURL(ctx, "https://google.com")

		Expect(err).ToNot(HaveOccurred())
	})

	Context("when the event cannot be stored", func() {
		It("should return the error", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			emitter.EXPECT().EmitShortURLCreated(gomock.Any(), "cv6VxVdu", "https://google.com").Return(errors.New("unknown error"))
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).To(MatchError("unable to save domain event: unknown error"))
			Expect(shortURL).To(BeNil())
		})
	})
})
