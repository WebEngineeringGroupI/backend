package redirect_test

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	eventmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

var _ = Describe("Domain / Redirect / LoadBalancerService", func() {
	var (
		ctx                   context.Context
		ctrl                  *gomock.Controller
		repository            *eventmocks.MockRepository
		multipleURLRedirector *redirect.LoadBalancerRedirectorService
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		repository = eventmocks.NewMockRepository(ctrl)
		multipleURLRedirector = redirect.NewLoadBalancerRedirectorService(repository)
		rand.Seed(GinkgoRandomSeed())
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	When("providing a hash", func() {
		It("returns only the valid URLs", func() {
			repository.EXPECT().
				Load(ctx, "someHash").
				Return(&url.LoadBalancedURL{Hash: "someHash", LongURLs: []url.OriginalURL{
					{
						URL:     "https://google.es",
						IsValid: true,
					},
					{
						URL:     "https://youtube.com",
						IsValid: false,
					},
				}}, 1, nil)

			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")

			Expect(err).ToNot(HaveOccurred())
			Expect(longURL).To(Equal("https://google.es"))
		})

		When("there are multiple valid URLs", func() {
			It("returns one of them randomly", func() {
				repository.EXPECT().
					Load(ctx, "someHash").
					AnyTimes().
					Return(&url.LoadBalancedURL{Hash: "someHash", LongURLs: []url.OriginalURL{
						{
							URL:     "https://google.es",
							IsValid: true,
						},
						{
							URL:     "https://youtube.com",
							IsValid: true,
						},
					}}, 3, nil)
				_, err := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")

				Expect(err).ToNot(HaveOccurred())
				Eventually(func() string {
					longURL, _ := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")
					return longURL
				}).Should(Equal("https://google.es"))
				Eventually(func() string {
					longURL, _ := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")
					return longURL
				}).Should(Equal("https://youtube.com"))
			})
		})
	})

	When("there are no valid URLs", func() {
		It("returns an error", func() {
			repository.EXPECT().
				Load(ctx, "someHash").
				Return(&url.LoadBalancedURL{Hash: "someHash", LongURLs: []url.OriginalURL{
					{
						URL:     "https://google.es",
						IsValid: false,
					},
					{
						URL:     "https://youtube.com",
						IsValid: false,
					},
				}}, 1, nil)
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")

			Expect(err).To(MatchError(url.ErrValidURLNotFound))
			Expect(longURL).To(BeEmpty())
		})
	})

	When("the repository returns an error", func() {
		It("returns the error from the repository", func() {
			repository.EXPECT().
				Load(ctx, "someHash").
				Return(nil, 0, fmt.Errorf("unknown error"))
			longURL, err := multipleURLRedirector.ReturnAValidOriginalURL(ctx, "someHash")

			Expect(err).To(MatchError("unknown error"))
			Expect(longURL).To(BeEmpty())
		})
	})
})
