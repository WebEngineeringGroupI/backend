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
	)

	BeforeEach(func() {
		repository = inmemory.NewRepository()
		validator = &FakeURLValidator{returnValidURL: true}
		shortener = url.NewSingleURLShortener(repository, validator)
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

		Context("and the provided URL is not valid", func() {
			It("validates that the provided URL is not valid", func() {
				aLongURL := "ftp://google.com"
				validator.shouldReturnValidURL(false)
				shortURL, err := shortener.HashFromURL(aLongURL)

				Expect(err).To(MatchError(url.ErrInvalidLongURLSpecified))
				Expect(shortURL).To(BeNil())
			})
		})

		Context("but the validator returns an error", func() {
			It("returns the error since it's unable to validate the URL", func() {
				aLongURL := "an-url"
				validator.shouldReturnError(errors.New("unknown error"))
				shortURL, err := shortener.HashFromURL(aLongURL)

				Expect(err).To(MatchError("unknown error"))
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

func (f *FakeURLValidator) ValidateURL(url string) (bool, error) {
	return f.returnValidURL, f.returnError
}

func (f *FakeURLValidator) ValidateURLs(url []string) (bool, error) {
	return f.returnValidURL, f.returnError
}
