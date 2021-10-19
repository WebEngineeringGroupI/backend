package redirect_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	`github.com/WebEngineeringGroupI/backend/pkg/domain/redirect`
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	`github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory`
)

var _ = Describe("Redirect", func() {
	var (
		repository url.ShortURLRepository
		redirector *redirect.Redirector
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		redirector = redirect.NewRedirector(repository)
	})

	Context("When providing a hash", func() {
		It("Returns a valid URL or an empty string", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "ftp://google.com",
			}
			repository.Save(shortURL)

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(originalURL).To(BeEmpty())
		})

		It("Returns a HTTP URL", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "http://google.com",
			}
			repository.Save(shortURL)

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)

			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal("http://google.com"))
		})

		It("Returns a HTTPS URL", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "https://google.com",
			}
			repository.Save(shortURL)

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal("https://google.com"))
		})
	})

	Context("When providing a hash", func() {
		It("Returns the same URL that generated the hash", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "https://google.com",
			}
			repository.Save(shortURL)

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)

			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal(shortURL.LongURL))
		})
	})

	Context("when providing a hash that doesn't exist", func() {
		It("the return value is an error", func() {
			_, err := redirector.ReturnOriginalURL("non-existing-hash")
			Expect(err).To(MatchError(url.ErrShortURLNotFound))
		})
	})
})
