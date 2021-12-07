package url_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Single URL shortener", func() {
	var (
		shortener  *url.SingleURLShortener
		repository url.ShortURLRepository
		validator  *FakeURLValidator
		metrics    *FakeMetrics
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		validator = &FakeURLValidator{returnValidURL: true}
		metrics = &FakeMetrics{}
		shortener = url.NewSingleURLShortener(repository, validator, metrics)
	})

	Context("when providing a long URL", func() {
		It("generates a hash", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.Hash).To(HaveLen(8))
			Expect(metrics.urlsProcessed).To(Equal(1))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		It("contains the real value from the original URL", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).ToNot(HaveOccurred())
			Expect(shortURL.OriginalURL.URL).To(Equal(aLongURL))
			Expect(metrics.urlsProcessed).To(Equal(1))
			Expect(metrics.singleURLMetrics).To(Equal(1))
		})

		Context("and the provided URL is not valid", func() {
			It("validates that the provided URL is not valid", func() {
				aLongURL := "ftp://google.com"
				validator.shouldReturnValidURL(false)
				shortURL, err := shortener.HashFromURL(aLongURL)

				Expect(err).To(MatchError(url.ErrInvalidLongURLSpecified))
				Expect(shortURL).To(BeNil())
				Expect(metrics.urlsProcessed).To(Equal(0))
				Expect(metrics.singleURLMetrics).To(Equal(1))
			})
		})

		Context("but the validator returns an error", func() {
			It("returns the error since it's unable to validate the URL", func() {
				aLongURL := "an-url"
				validator.shouldReturnError(errors.New("unknown error"))
				shortURL, err := shortener.HashFromURL(aLongURL)

				Expect(err).To(MatchError("unknown error"))
				Expect(shortURL).To(BeNil())
				Expect(metrics.urlsProcessed).To(Equal(0))
				Expect(metrics.singleURLMetrics).To(Equal(1))
			})
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				shortGoogleURL, err := shortener.HashFromURL("https://google.com")
				Expect(err).ToNot(HaveOccurred())

				shortFacebookURL, err := shortener.HashFromURL("https://facebook.com")
				Expect(err).ToNot(HaveOccurred())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
				Expect(metrics.urlsProcessed).To(Equal(2))
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

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})
})
