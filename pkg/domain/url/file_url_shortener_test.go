package url_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Multiple URL Shortener", func() {
	var (
		shortener  *url.FileURLShortener
		repository url.ShortURLRepository
		formatter  *FakeFormatter
		metrics    *FakeMetrics
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		formatter = &FakeFormatter{}
		metrics = &FakeMetrics{}
		formatter.shouldReturnURLs(aLongURLSet())
		shortener = url.NewFileURLShortener(repository, metrics, formatter)
	})

	Context("when providing multiple long URLs", func() {
		It("generates a hash for each one", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs).To(HaveLen(2))
			Expect(metrics.fileURLMetrics).To(Equal(1))
			Expect(shortURLs[0].Hash).To(HaveLen(8))
			Expect(shortURLs[1].Hash).To(HaveLen(8))
		})

		It("contains the real values from the original URLs", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(metrics.fileURLMetrics).To(Equal(1))
			Expect(shortURLs[0].OriginalURL.URL).To(Equal("https://google.com"))
			Expect(shortURLs[1].OriginalURL.URL).To(Equal("https://unizar.es"))
		})

		It("saves the URLs as not verified", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURLs[0].OriginalURL.IsValid).To(BeFalse())
			Expect(shortURLs[1].OriginalURL.IsValid).To(BeFalse())
		})

		It("generates different short URL hashes for each of the long URLs", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())

			Expect(err).ToNot(HaveOccurred())
			Expect(metrics.fileURLMetrics).To(Equal(1))
			Expect(shortURLs[0].Hash).ToNot(Equal(shortURLs[1].Hash))
		})

		Context("but the provided data is not valid", func() {
			It("returns the error since it's unable to transform the data", func() {
				formatter.shouldReturnError(errors.New("unknown error"))
				shortURLs, err := shortener.HashesFromURLData(aLongURLData())

				Expect(err).To(MatchError("unknown error"))
				Expect(shortURLs).To(BeNil())
				Expect(metrics.fileURLMetrics).To(Equal(1))
			})
		})

		It("stores the short URL in a repository", func() {
			shortURLs, err := shortener.HashesFromURLData(aLongURLData())
			Expect(err).ToNot(HaveOccurred())

			firstURL, err := repository.FindShortURLByHash(shortURLs[0].Hash)
			Expect(err).ToNot(HaveOccurred())
			secondURL, err := repository.FindShortURLByHash(shortURLs[1].Hash)
			Expect(err).ToNot(HaveOccurred())

			Expect(firstURL.Hash).To(Equal(shortURLs[0].Hash))
			Expect(firstURL.OriginalURL.URL).To(Equal("https://google.com"))
			Expect(secondURL.Hash).To(Equal(shortURLs[1].Hash))
			Expect(secondURL.OriginalURL.URL).To(Equal("https://unizar.es"))
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
