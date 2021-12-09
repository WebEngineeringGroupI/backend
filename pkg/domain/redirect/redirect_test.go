package redirect_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Redirect", func() {
	var (
		repository url.ShortURLRepository
		redirector *redirect.Redirector
		ctx        context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		repository = inmemory.NewRepository()
		redirector = redirect.NewRedirector(repository)
	})

	Context("when providing a hash", func() {
		It("returns a HTTP URL", func() {
			shortURL := &url.ShortURL{
				Hash:        "asdfasdf",
				OriginalURL: url.OriginalURL{URL: "http://google.com", IsValid: true},
			}
			err := repository.SaveShortURL(ctx, shortURL)
			Expect(err).ToNot(HaveOccurred())

			originalURL, err := redirector.ReturnOriginalURL(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(originalURL).To(Equal("http://google.com"))
		})

		It("returns a HTTPS URL", func() {
			shortURL := &url.ShortURL{
				Hash:        "asdfasdf",
				OriginalURL: url.OriginalURL{URL: "https://google.com", IsValid: true},
			}

			err := repository.SaveShortURL(ctx, shortURL)
			Expect(err).ToNot(HaveOccurred())

			originalURL, err := redirector.ReturnOriginalURL(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(originalURL).To(Equal("https://google.com"))
		})
	})

	Context("when providing a hash", func() {
		It("returns the same URL that generated the hash", func() {
			shortURL := &url.ShortURL{
				Hash:        "asdfasdf",
				OriginalURL: url.OriginalURL{URL: "http://google.com", IsValid: true},
			}

			err := repository.SaveShortURL(ctx, shortURL)
			Expect(err).ToNot(HaveOccurred())

			originalURL, err := redirector.ReturnOriginalURL(ctx, shortURL.Hash)
			Expect(err).ToNot(HaveOccurred())
			Expect(originalURL).To(Equal(shortURL.OriginalURL.URL))
		})
	})

	Context("if the URL is not valid", func() {
		It("returns an error saying it's not valid", func() {
			shortURL := &url.ShortURL{
				Hash:        "12345",
				OriginalURL: url.OriginalURL{URL: "some-url", IsValid: false},
			}
			_ = repository.SaveShortURL(ctx, shortURL)

			originalURL, err := redirector.ReturnOriginalURL(ctx, "12345")

			Expect(err).To(MatchError("the url 'some-url' is marked as invalid"))
			Expect(originalURL).To(BeEmpty())
		})
	})

	Context("when providing a hash that doesn't exist", func() {
		It("the return value is an error", func() {
			_, err := redirector.ReturnOriginalURL(ctx, "non-existing-hash")
			Expect(err).To(MatchError(url.ErrShortURLNotFound))
		})
	})
})
