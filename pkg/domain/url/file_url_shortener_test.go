package url_test

import (
	"context"
	"errors"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	domainmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
)

var _ = Describe("Multiple URL Shortener", func() {
	var (
		ctrl       *gomock.Controller
		shortener  *url.FileURLShortener
		repository *domainmocks.MockRepository
		formatter  *mocks.MockFormatter
		clock      *domainmocks.MockClock
		metrics    *mocks.MockMetrics
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		metrics = mocks.NewMockMetrics(ctrl)
		formatter = mocks.NewMockFormatter(ctrl)
		clock = domainmocks.NewMockClock(ctrl)
		repository = domainmocks.NewMockRepository(ctrl)

		shortener = url.NewFileURLShortener(repository, metrics, clock, formatter)

		clock.EXPECT().Now().AnyTimes().Return(time.Time{})
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when providing multiple long URLs", func() {
		BeforeEach(func() {
			formatter.EXPECT().FormatDataToURLs(gomock.Any()).Return(aLongURLSet(), nil)
		})

		It("generates a hash for each one", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any()).AnyTimes()
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs).To(HaveLen(2))
			Expect(shortURLs[0].Hash).To(HaveLen(8))
			Expect(shortURLs[1].Hash).To(HaveLen(8))
		})

		It("contains the real values from the original URLs", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any()).AnyTimes()
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].OriginalURL.URL).To(Equal("https://google.com"))
			Expect(shortURLs[1].OriginalURL.URL).To(Equal("https://unizar.es"))
		})

		It("saves the URLs as not verified", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any()).AnyTimes()
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].OriginalURL.IsValid).To(BeFalse())
			Expect(shortURLs[1].OriginalURL.IsValid).To(BeFalse())
		})

		It("generates different short URL hashes for each of the long URLs", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, gomock.Any()).AnyTimes()
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].Hash).ToNot(Equal(shortURLs[1].Hash))
		})

		It("stores the short URL in a repository", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			repository.EXPECT().Save(ctx, []event.Event{
				&url.ShortURLCreated{
					Base: event.Base{
						ID:      "cv6VxVdu",
						Version: 0,
						At:      time.Time{},
					},
					OriginalURL: "https://google.com",
				},
			})
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

			_, err := shortener.HashesFromURLData(ctx, aLongURLData())
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("when the provided data is not valid", func() {
		BeforeEach(func() {
			formatter.EXPECT().FormatDataToURLs(gomock.Any()).Return(nil, errors.New("unknown error"))
		})
		It("returns the error since it's unable to transform the data", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).To(MatchError("unknown error"))
			Expect(shortURLs).To(BeNil())
		})
	})
})

func aLongURLData() []byte {
	return []byte(`"https://google.com"
"https://unizar.es"`)
}

func aLongURLSet() []string {
	return []string{"https://google.com", "https://unizar.es"}
}
