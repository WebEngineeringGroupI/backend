package redirect_test

import (
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type FakeMultipleShortURLsRepository struct {
	validURLs     []string
	invalidURLs   []string
	errorToReturn error
}

var _ = Describe("Domain / Redirect / LoadBalancer", func() {
	var (
		repository            *FakeMultipleShortURLsRepository
		multipleURLRedirector *redirect.LoadBalancerRedirector
	)

	BeforeEach(func() {
		repository = &FakeMultipleShortURLsRepository{}
		multipleURLRedirector = redirect.NewLoadBalancerRedirector(repository)
		rand.Seed(GinkgoRandomSeed())
	})

	When("providing a hash", func() {
		It("returns only the valid URLs", func() {
			repository.shouldReturnValidURLs("https://google.es").shouldReturnInvalidURLs("https://youtube.com")
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL("someHash")

			Expect(err).ToNot(HaveOccurred())
			Expect(longURL).To(Equal("https://google.es"))
		})

		When("there are multiple valid URLs", func() {
			It("returns one of them randomly", func() {
				repository.shouldReturnValidURLs("https://google.es", "https://youtube.com")
				_, err := multipleURLRedirector.ReturnAValidOriginalURL("someHash")

				Expect(err).ToNot(HaveOccurred())
				Eventually(func() string { longURL, _ := multipleURLRedirector.ReturnAValidOriginalURL("someHash"); return longURL }).Should(Equal("https://google.es"))
				Eventually(func() string { longURL, _ := multipleURLRedirector.ReturnAValidOriginalURL("someHash"); return longURL }).Should(Equal("https://youtube.com"))
			})
		})
	})

	When("there are no valid URLs", func() {
		It("returns an error", func() {
			repository.shouldReturnInvalidURLs("https://youtube.com", "https://google.es")
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL("someHash")

			Expect(err).To(MatchError("there are no valid URLs to redirect to"))
			Expect(longURL).To(BeEmpty())
		})
	})

	When("the repository does not find a valid URL for the hash", func() {
		It("returns the error", func() {
			repository.shouldReturnError(redirect.ErrValidURLNotFound)
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL("someHash")

			Expect(err).To(MatchError(redirect.ErrValidURLNotFound))
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

func (f *FakeMultipleShortURLsRepository) FindByHash(hash string) (*url.LoadBalancedURL, error) {
	if f.errorToReturn != nil {
		return nil, f.errorToReturn
	}

	originalURLs := []url.OriginalURL{}
	for _, aURL := range f.validURLs {
		originalURLs = append(originalURLs, url.OriginalURL{
			URL:     aURL,
			IsValid: true,
		})
	}
	for _, aURL := range f.invalidURLs {
		originalURLs = append(originalURLs, url.OriginalURL{
			URL:     aURL,
			IsValid: false,
		})
	}

	return &url.LoadBalancedURL{Hash: hash, LongURLs: originalURLs}, nil
}

func (f *FakeMultipleShortURLsRepository) Save(urls *url.LoadBalancedURL) error {
	panic("implement me")
}
