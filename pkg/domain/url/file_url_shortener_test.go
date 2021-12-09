package url_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Multiple URL Shortener", func() {
	var (
		controller *gomock.Controller
		shortener  *url.FileURLShortener
		repository url.ShortURLRepository
		formatter  *mocks.MockFormatter
		metrics    *mocks.MockMetrics
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repository = inmemory.NewRepository()

		controller = gomock.NewController(GinkgoT())
		metrics = mocks.NewMockMetrics(controller)
		formatter = mocks.NewMockFormatter(controller)

		shortener = url.NewFileURLShortener(repository, metrics, formatter)
	})
	AfterEach(func() {
		controller.Finish()
	})

	Context("when providing multiple long URLs", func() {
		BeforeEach(func() {
			formatter.EXPECT().FormatDataToURLs(gomock.Any()).Return(aLongURLSet(), nil)
		})

		It("generates a hash for each one", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs).To(HaveLen(2))
			Expect(shortURLs[0].Hash).To(HaveLen(8))
			Expect(shortURLs[1].Hash).To(HaveLen(8))
		})

		It("contains the real values from the original URLs", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].OriginalURL.URL).To(Equal("https://google.com"))
			Expect(shortURLs[1].OriginalURL.URL).To(Equal("https://unizar.es"))
		})

		It("saves the URLs as not verified", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].OriginalURL.IsValid).To(BeFalse())
			Expect(shortURLs[1].OriginalURL.IsValid).To(BeFalse())
		})

		It("generates different short URL hashes for each of the long URLs", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].Hash).ToNot(Equal(shortURLs[1].Hash))
		})

		It("stores the short URL in a repository", func() {
			metrics.EXPECT().RecordFileURLMetrics().Times(1)
			shortURLs, err := shortener.HashesFromURLData(ctx, aLongURLData())
			Expect(err).ToNot(HaveOccurred())

			firstURL, err := repository.FindShortURLByHash(ctx, shortURLs[0].Hash)
			Expect(err).ToNot(HaveOccurred())
			secondURL, err := repository.FindShortURLByHash(ctx, shortURLs[1].Hash)
			Expect(err).ToNot(HaveOccurred())

			Expect(firstURL.Hash).To(Equal(shortURLs[0].Hash))
			Expect(firstURL.OriginalURL.URL).To(Equal("https://google.com"))
			Expect(secondURL.Hash).To(Equal(shortURLs[1].Hash))
			Expect(secondURL.OriginalURL.URL).To(Equal("https://unizar.es"))
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
