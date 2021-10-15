package redirect_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redirect", func() {

	repository := &FakeShortURLRepository{ urls: map[string]*url.ShortURL{} }
	redirector := redirect.NewRedirector(repository)

	Context("When providing a hash", func() {
		It("Returns a valid URL or an empty string", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash: "asdfasdf",
				LongURL: "ftp://google.com",
			}
			repository.Save(shortURL)

			Expect(redirector.ReturnOriginalURL(shortURL.Hash)).To(BeEmpty())
		})
		It("Returns a HTTP URL", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash: "asdfasdf",
				LongURL: "http://google.com",
			}
			repository.Save(shortURL)

			Expect(redirector.ReturnOriginalURL(shortURL.Hash)).To(Equal("http://google.com"))
		})
		It("Returns a HTTP URL", func() {

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash: "asdfasdf",
				LongURL: "https://google.com",
			}
			repository.Save(shortURL)

			Expect(redirector.ReturnOriginalURL(shortURL.Hash)).To(Equal("https://google.com"))
		})
	})

	Context("When providing a hash", func() {
		It("Returns the same URL that generated the hash", func(){

			// Generate a Hash and a shortURL from a URL
			shortURL := &url.ShortURL{
				Hash: "asdfasdf",
				LongURL: "https://google.com",
			}
			repository.Save(shortURL)

			Expect(redirector.ReturnOriginalURL(shortURL.Hash)).To(Equal(shortURL.LongURL))
		})
	})
})

type FakeShortURLRepository struct {
	urls map[string]*url.ShortURL
}

func (f *FakeShortURLRepository) Save(url *url.ShortURL) {
	f.urls[url.Hash] = url
}

func (f *FakeShortURLRepository) FindByHash(hash string) *url.ShortURL {
	url, ok := f.urls[hash]
	if !ok {
		return nil
	}

	return url
}