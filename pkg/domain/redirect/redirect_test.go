package redirect_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Redirect", func() {
	var (
		repository url.ShortURLRepository
		validator  *FakeURLValidator
		redirector *redirect.Redirector
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		validator = &FakeURLValidator{returnValidURL: true}
		redirector = redirect.NewRedirector(repository, validator)
	})

	Context("when providing a hash", func() {
		It("returns a HTTP URL", func() {
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "http://google.com",
			}
			err := repository.Save(shortURL)
			Expect(err).To(Succeed())

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal("http://google.com"))
		})

		It("returns a HTTPS URL", func() {
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "https://google.com",
			}

			err := repository.Save(shortURL)
			Expect(err).To(Succeed())

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal("https://google.com"))
		})
	})

	Context("when providing a hash", func() {
		It("returns the same URL that generated the hash", func() {
			shortURL := &url.ShortURL{
				Hash:    "asdfasdf",
				LongURL: "https://google.com",
			}

			err := repository.Save(shortURL)
			Expect(err).To(Succeed())

			originalURL, err := redirector.ReturnOriginalURL(shortURL.Hash)
			Expect(err).To(Succeed())
			Expect(originalURL).To(Equal(shortURL.LongURL))
		})
	})

	Context("when validating the URL", func() {
		Context("if the URL is not valid", func() {
			It("returns an error saying it's not valid", func() {
				shortURL := &url.ShortURL{
					Hash:    "12345",
					LongURL: "some-url",
				}
				_ = repository.Save(shortURL)
				validator.shouldReturnValidURL(false)

				originalURL, err := redirector.ReturnOriginalURL("12345")

				Expect(err).To(MatchError("the url 'some-url' is marked as invalid"))
				Expect(originalURL).To(BeEmpty())
			})
		})
		Context("if the validator is not able to validate the URL", func() {
			It("returns the error saying it's not able to validate it", func() {
				shortURL := &url.ShortURL{
					Hash:    "12345",
					LongURL: "some-url",
				}
				_ = repository.Save(shortURL)
				validator.shouldReturnError(errors.New("unknown validation error"))

				originalURL, err := redirector.ReturnOriginalURL("12345")

				Expect(err).To(MatchError("unknown validation error"))
				Expect(originalURL).To(BeEmpty())
			})
		})
	})

	Context("when providing a hash that doesn't exist", func() {
		It("the return value is an error", func() {
			_, err := redirector.ReturnOriginalURL("non-existing-hash")
			Expect(err).To(MatchError(url.ErrShortURLNotFound))
		})
	})
})

type FakeURLValidator struct {
	returnValidURL bool
	returnError    error
}

func (f *FakeURLValidator) shouldReturnValidURL(validURL bool) {
	f.returnValidURL = validURL
}

func (f *FakeURLValidator) shouldReturnError(err error) {
	f.returnError = err
}

func (f *FakeURLValidator) ValidateURLs(urls []string) (bool, error) {
	return f.returnValidURL, f.returnError
}
