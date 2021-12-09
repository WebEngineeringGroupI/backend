package url_test

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
)

var _ = Describe("Domain / URL / Load Balancing", func() {
	var (
		loadBalancer                *url.LoadBalancer
		controller                  *gomock.Controller
		multipleShortURLsRepository *mocks.MockLoadBalancedURLsRepository
		ctx                         context.Context
	)
	BeforeEach(func() {
		controller = gomock.NewController(GinkgoT())
		multipleShortURLsRepository = mocks.NewMockLoadBalancedURLsRepository(controller)
		loadBalancer = url.NewLoadBalancer(multipleShortURLsRepository)
		ctx = context.Background()
	})
	AfterEach(func() {
		controller.Finish()
	})

	When("a single URL is generated from multiple URLs", func() {
		It("is correctly generated", func() {
			multipleShortURLsRepository.EXPECT().SaveLoadBalancedURL(gomock.Any(), &url.LoadBalancedURL{
				Hash: "P3Z83Gpy",
				LongURLs: []url.OriginalURL{
					{URL: "aURL", IsValid: false},
					{URL: "anotherURL", IsValid: false},
				},
			})

			loadBalancedURLs, err := loadBalancer.ShortURLs(ctx, []string{"aURL", "anotherURL"})

			Expect(err).ToNot(HaveOccurred())
			Expect(loadBalancedURLs).To(Equal(&url.LoadBalancedURL{
				Hash: "P3Z83Gpy",
				LongURLs: []url.OriginalURL{
					{URL: "aURL", IsValid: false},
					{URL: "anotherURL", IsValid: false},
				},
			}))
		})
	})

	When("the repository returns an error", func() {
		It("returns the error from the repository", func() {
			multipleShortURLsRepository.EXPECT().SaveLoadBalancedURL(gomock.Any(), gomock.Any()).Return(errors.New("unknown error"))
			loadBalancedURLs, err := loadBalancer.ShortURLs(ctx, []string{"aURL"})

			Expect(err).To(MatchError("error saving load-balanced URLs into repository: unknown error"))
			Expect(loadBalancedURLs).To(BeNil())
		})
	})

	When("the list of URLs is empty", func() {
		It("returns an error", func() {
			loadBalancedURLs, err := loadBalancer.ShortURLs(ctx, []string{})

			Expect(err).To(MatchError(url.ErrNoURLsSpecified))
			Expect(loadBalancedURLs).To(BeNil())
		})
	})

	When("the list has more than the allowed number of elements", func() {
		It("returns an error saying no so many elements are allowed", func() {
			multipleShortURLs, err := loadBalancer.ShortURLs(ctx, urlListOfSize(11))

			Expect(err).To(MatchError(url.ErrTooMuchMultipleURLs))
			Expect(multipleShortURLs).To(BeNil())
		})
	})
})

func urlListOfSize(size int) []string {
	list := make([]string, 0, size)
	for i := 0; i < size; i++ {
		list = append(list, fmt.Sprintf("anURL%d", i))
	}
	return list
}
