package url_test

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	domainmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Single URL shortener", func() {
	var (
		controller *gomock.Controller
		shortener  *url.SingleURLShortener
		repository url.ShortURLRepository
		metrics    *FakeMetrics
		clock      *mocks.MockClock
		uuid       *mocks.MockUUID
		outbox     *domainmocks.MockEventOutbox
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		controller = gomock.NewController(GinkgoT())
		clock = mocks.NewMockClock(controller)
		uuid = mocks.NewMockUUID(controller)
		outbox = domainmocks.NewMockEventOutbox(controller)
		repository = inmemory.NewRepository()
		metrics = &FakeMetrics{}
		shortener = url.NewSingleURLShortenerToTest(repository, metrics, outbox, clock, uuid)

		clock.EXPECT().Now().Return(time.Time{}).AnyTimes()
		uuid.EXPECT().New().Return("aUUID").AnyTimes()
	})

	AfterEach(func() {
		controller.Finish()
	})

	Context("when providing a long URL", func() {
		BeforeEach(func() {
			outbox.EXPECT().SaveEvent(gomock.Any(), gomock.Any()).AnyTimes()
		})

		It("generates a hash", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(ctx, aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.Hash).To(Equal("cv6VxVdu"))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		It("contains the real value from the original URL", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(ctx, aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.OriginalURL.URL).To(Equal(aLongURL))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				shortGoogleURL, err := shortener.HashFromURL(ctx, "https://google.com")
				Expect(err).ToNot(HaveOccurred())

				shortFacebookURL, err := shortener.HashFromURL(ctx, "https://facebook.com")
				Expect(err).ToNot(HaveOccurred())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
				Expect(metrics.singleURLMetrics).To(Equal(2))
			})
		})

		It("stores the short URL in a repository", func() {
			shortURL, err := shortener.HashFromURL(ctx, "https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.Hash).To(Equal(shortURL.Hash))
		})

		It("stores the URL as non verified", func() {
			shortURL, err := shortener.HashFromURL(ctx, "https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.OriginalURL.IsValid).To(BeFalse())
		})

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})

	It("should emit an event saved to the outbox", func() {
		outbox.EXPECT().SaveEvent(gomock.Any(), url.NewShortURLCreated(&url.ShortURL{
			Hash: "cv6VxVdu",
			OriginalURL: url.OriginalURL{
				URL:     "https://google.com",
				IsValid: false,
			},
		}))
		_, err := shortener.HashFromURL(ctx, "https://google.com")

		Expect(err).ToNot(HaveOccurred())
	})

	Context("when the event cannot be stored", func() {
		It("should return the error", func() {
			outbox.EXPECT().SaveEvent(gomock.Any(), gomock.Any()).Return(errors.New("unknown error"))
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).To(MatchError("unable to save domain event: unknown error"))
			Expect(shortURL).To(BeNil())
		})
	})
})
