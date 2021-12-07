package redirect_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"math/rand"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
)

type FakeMultipleShortURLsRepository struct {
	validURLs     []string
	invalidURLs   []string
	errorToReturn error
}

var _ = Describe("Domain / Redirect / LoadBalancer", func() {
	var (
		repository            *FakeMultipleShortURLsRepository
		multipleURLRedirector *redirect.MultipleURLRedirector
	)

	BeforeEach(func() {
		repository = &FakeMultipleShortURLsRepository{}
		multipleURLRedirector = redirect.NewMultipleURLRedirector(repository)
		rand.Seed(GinkgoRandomSeed())
	})

	When("providing a hash", func() {
		It("returns only the valid URLs", func() {
			repository.shouldReturnValidURLs("https://google.es").shouldReturnInvalidURLs("https://youtube.com")
			longURL, err := multipleURLRedirector.ReturnValidOriginalURL("someHash")

			Expect(err).ToNot(HaveOccurred())
			Expect(longURL).To(Equal("https://google.es"))
		})

		When("there are multiple valid URLs", func() {
			It("returns one of them randomly", func() {
				repository.shouldReturnValidURLs("https://google.es", "https://youtube.com")
				_, err := multipleURLRedirector.ReturnValidOriginalURL("someHash")

				Expect(err).ToNot(HaveOccurred())
				Eventually(func() string { longURL, _ := multipleURLRedirector.ReturnValidOriginalURL("someHash"); return longURL }).Should(Equal("https://google.es"))
				Eventually(func() string { longURL, _ := multipleURLRedirector.ReturnValidOriginalURL("someHash"); return longURL }).Should(Equal("https://youtube.com"))
			})
		})
	})

	When("there are no valid URLs", func() {
		It("returns an error", func() {
			repository.shouldReturnInvalidURLs("https://youtube.com", "https://google.es")
			longURL, err := multipleURLRedirector.ReturnValidOriginalURL("someHash")

			Expect(err).To(MatchError("there are no valid URLs to redirect to"))
			Expect(longURL).To(BeEmpty())
		})
	})

	When("the repository does not find a valid URL for the hash", func() {
		It("returns the error", func() {
			repository.shouldReturnError(redirect.ValidURLNotFound)
			longURL, err := multipleURLRedirector.ReturnValidOriginalURL("someHash")

			Expect(err).To(MatchError(redirect.ValidURLNotFound))
			Expect(longURL).To(BeEmpty())
		})
	})
})

func (f *FakeMultipleShortURLsRepository) shouldReturnValidURLs(urls ...string) *FakeMultipleShortURLsRepository {
	f.validURLs = append(f.validURLs, urls...)
	return f
}
func (f *FakeMultipleShortURLsRepository) shouldReturnInvalidURLs(urls ...string) *FakeMultipleShortURLsRepository {
	f.invalidURLs = append(f.invalidURLs, urls...)
	return f
}

func (f *FakeMultipleShortURLsRepository) shouldReturnError(err error) *FakeMultipleShortURLsRepository {
	f.errorToReturn = err
	return f
}

func (f *FakeMultipleShortURLsRepository) FindOriginalURLsForHash(hash string) ([]url.OriginalURL, error) {
	if f.errorToReturn != nil {
		return nil, f.errorToReturn
	}

	result := []url.OriginalURL{}
	for _, aURL := range f.validURLs {
		result = append(result, url.OriginalURL{
			URL:     aURL,
			IsValid: true,
		})
	}
	for _, aURL := range f.invalidURLs {
		result = append(result, url.OriginalURL{
			URL:     aURL,
			IsValid: false,
		})
	}

	return result, nil
}
