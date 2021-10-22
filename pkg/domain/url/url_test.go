package url_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("URL shortener", func() {
	var (
		shortener  *url.Shortener
		repository url.ShortURLRepository
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		shortener = url.NewShortener(repository)
	})

	Context("when providing a long URL", func() {
		It("generates a hash", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).To(Succeed())
			Expect(shortURL.Hash).To(HaveLen(8))
		})

		It("contains the real value from the original URL", func() {
			aLongURL := "https://google.com"
			shortURL, err := shortener.HashFromURL(aLongURL)

			Expect(err).To(Succeed())
			Expect(shortURL.LongURL).To(Equal(aLongURL))
		})

		Context("and the provided URL is not HTTP", func() {
			It("validates that the provided URL is not valid", func() {
				aLongURL := "ftp://google.com"
				shortURL, err := shortener.HashFromURL(aLongURL)

				Expect(err).To(MatchError(url.ErrInvalidLongURLSpecified))
				Expect(shortURL).To(BeNil())
			})
		})

		Context("when providing different long URLs", func() {
			It("generates different short URL hashes", func() {
				shortGoogleURL, err := shortener.HashFromURL("https://google.com")
				Expect(err).To(Succeed())

				shortFacebookURL, err := shortener.HashFromURL("https://facebook.com")
				Expect(err).To(Succeed())

				Expect(shortGoogleURL.Hash).ToNot(Equal(shortFacebookURL.Hash))
			})
		})

		It("stores the short URL in a repository", func() {
			shortURL, err := shortener.HashFromURL("https://unizar.es")
			Expect(err).To(Succeed())

			expectedURLInRepo, err := repository.FindByHash(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(expectedURLInRepo.Hash).To(Equal(shortURL.Hash))
		})

		// TODO(german): Each time a new hash is generated, do we need to check if it already exists?
		// TODO(german): What's the meaning of Safe and Sponsor in the original urlshortener implementation
	})
})
