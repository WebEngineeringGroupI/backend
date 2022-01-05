package url_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	domainmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
)

var _ = Describe("Single URL shortener", func() {
	var (
		ctrl       *gomock.Controller
		shortener  *url.SingleURLShortener
		repository *domainmocks.MockRepository
		clock      *domainmocks.MockClock
		metrics    *mocks.MockMetrics
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		repository = domainmocks.NewMockRepository(ctrl)
		clock = domainmocks.NewMockClock(ctrl)
		metrics = mocks.NewMockMetrics(ctrl)

		shortener = url.NewSingleURLShortener(repository, clock, metrics)

		clock.EXPECT().Now().AnyTimes().Return(time.Time{})
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when providing a long URL", func() {

		It("generates a hash", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any())
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.Hash).To(Equal("cv6VxVdu"))
		})

		It("contains the real value from the original URL", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any())
			shortURL, err := shortener.HashFromURL(ctx, "https://google.com")

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.OriginalURL.URL).To(Equal("https://google.com"))
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				metrics.EXPECT().RecordSingleURLMetrics().Times(2)
				repository.EXPECT().Save(ctx, gomock.Any()).Times(2)
				shortGoogleURL, err := shortener.HashFromURL(ctx, "https://google.com")
				Expect(err).ToNot(HaveOccurred())

				shortFacebookURL, err := shortener.HashFromURL(ctx, "https://facebook.com")
				Expect(err).ToNot(HaveOccurred())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
			})
		})

		It("stores the short URL in a repository", func() {
			metrics.EXPECT().RecordSingleURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, []event.Event{
				&url.ShortURLCreated{
					Base: event.Base{
						ID:      "2sMi6l0Z",
						Version: 0,
						At:      time.Time{},
					},
					OriginalURL: "https://unizar.es",
				},
			})
			shortURL, err := shortener.HashFromURL(ctx, "https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			Expect(shortURL.Hash).To(Equal("2sMi6l0Z"))
		})

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})
})
