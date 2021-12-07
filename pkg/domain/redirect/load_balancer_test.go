package redirect_test

import (
	"errors"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type FakeLoadBalancedURLsRepository struct {
	validURLs     []string
	invalidURLs   []string
	errorToReturn error
}

var _ = Describe("Domain / Redirect / LoadBalancer", func() {
	var (
		repository            *FakeLoadBalancedURLsRepository
		multipleURLRedirector *redirect.LoadBalancerRedirector
	)

	BeforeEach(func() {
		repository = &FakeLoadBalancedURLsRepository{}
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

			Expect(err).To(MatchError(url.ErrValidURLNotFound))
			Expect(longURL).To(BeEmpty())
		})
	})

	When("the repository returns an error", func() {
		It("returns the error from the repository", func() {
			repository.shouldReturnError(errors.New("unknown error"))
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL("someHash")

			Expect(err).To(MatchError("unknown error"))
			Expect(longURL).To(BeEmpty())
		})
	})
})

func (f *FakeLoadBalancedURLsRepository) shouldReturnValidURLs(urls ...string) *FakeLoadBalancedURLsRepository {
	f.validURLs = append(f.validURLs, urls...)
	return f
}

func (f *FakeLoadBalancedURLsRepository) shouldReturnInvalidURLs(urls ...string) *FakeLoadBalancedURLsRepository {
	f.invalidURLs = append(f.invalidURLs, urls...)
	return f
}
func (f *FakeLoadBalancedURLsRepository) shouldReturnError(err error) *FakeLoadBalancedURLsRepository {
	f.errorToReturn = err
	return f
}

func (f *FakeLoadBalancedURLsRepository) FindByHash(hash string) (*url.LoadBalancedURL, error) {
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

func (f *FakeLoadBalancedURLsRepository) Save(urls *url.LoadBalancedURL) error {
	panic("implement me")
}
