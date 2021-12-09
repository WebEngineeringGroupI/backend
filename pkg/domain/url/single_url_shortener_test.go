package url_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Single URL shortener", func() {
	var (
		shortener  *url.SingleURLShortener
		repository url.ShortURLRepository
		metrics    *FakeMetrics
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		metrics = &FakeMetrics{}
		shortener = url.NewSingleURLShortener(repository, metrics)
	})

	Context("when providing a long URL", func() {
		It("generates a hash", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.Hash).To(HaveLen(8))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		It("contains the real value from the original URL", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.OriginalURL.URL).To(Equal(aLongURL))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				shortGoogleURL, err := shortener.HashFromURL("https://google.com")
				Expect(err).ToNot(HaveOccurred())

				shortFacebookURL, err := shortener.HashFromURL("https://facebook.com")
				Expect(err).ToNot(HaveOccurred())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
				Expect(metrics.singleURLMetrics).To(Equal(2))
			})
		})

		It("stores the short URL in a repository", func() {
			shortURL, err := shortener.HashFromURL("https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.Hash).To(Equal(shortURL.Hash))
		})

		It("stores the URL as non verified", func() {
			shortURL, err := shortener.HashFromURL("https://unizar.es")
			Expect(err).ToNot(HaveOccurred())

			expectedURLInRepo, err := repository.FindShortURLByHash(shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(expectedURLInRepo.OriginalURL.IsValid).To(BeFalse())
		})

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})
})
