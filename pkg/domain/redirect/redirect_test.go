package redirect_test

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	domainmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

var _ = Describe("Redirect", func() {
	var (
		ctrl       *gomock.Controller
		repository *domainmocks.MockRepository
		clock      *domainmocks.MockClock
		redirector *redirect.Redirector
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		repository = domainmocks.NewMockRepository(ctrl)
		clock = domainmocks.NewMockClock(ctrl)

		redirector = redirect.NewRedirector(repository, clock)

		clock.EXPECT().Now().AnyTimes().Return(time.Time{})
	})

	Context("when providing a hash", func() {
		It("returns a HTTP URL", func() {
			repository.EXPECT().Load(ctx, "foobar").Return(&url.ShortURL{
				Hash: "foobar",
				OriginalURL: url.OriginalURL{
					URL:     "http://google.com",
					IsValid: true,
				},
				Clicks: 2,
			}, 3, nil)
			repository.EXPECT().Save(ctx, &url.ShortURLClicked{
				Base: event.Base{
					ID:      "foobar",
					Version: 4,
					At:      time.Time{},
				},
			})

			originalURL, err := redirector.ReturnOriginalURL(ctx, "foobar")

			Expect(err).ToNot(HaveOccurred())
			Expect(originalURL).To(Equal("http://google.com"))
		})

		It("returns a HTTPS URL", func() {
			repository.EXPECT().Load(ctx, "foobar").Return(&url.ShortURL{
				Hash: "foobar",
				OriginalURL: url.OriginalURL{
					URL:     "https://google.com",
					IsValid: true,
				},
				Clicks: 1,
			}, 6, nil)
			repository.EXPECT().Save(ctx, &url.ShortURLClicked{
				Base: event.Base{
					ID:      "foobar",
					Version: 7,
					At:      time.Time{},
				},
			})

			originalURL, err := redirector.ReturnOriginalURL(ctx, "foobar")

			Expect(err).ToNot(HaveOccurred())
			Expect(originalURL).To(Equal("https://google.com"))
		})
	})

	Context("if the URL is not valid", func() {
		It("returns an error saying it's not valid", func() {
			repository.EXPECT().Load(ctx, "12345").Return(&url.ShortURL{
				Hash: "12345",
				OriginalURL: url.OriginalURL{
					URL:     "some-url",
					IsValid: false,
				},
				Clicks: 1,
			}, 6, nil)

			originalURL, err := redirector.ReturnOriginalURL(ctx, "12345")

			Expect(err).To(MatchError("the url 'some-url' is marked as invalid"))
			Expect(originalURL).To(BeEmpty())
		})
	})

	Context("when providing a hash that doesn't exist", func() {
		It("the return value is an error", func() {
			repository.EXPECT().Load(ctx, "non-existing-hash").Return(nil, 0, url.ErrShortURLNotFound)

			_, err := redirector.ReturnOriginalURL(ctx, "non-existing-hash")
			Expect(err).To(MatchError(url.ErrShortURLNotFound))
		})
	})
})
