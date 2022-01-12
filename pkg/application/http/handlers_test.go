package http_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"math/rand"
	gohttp "net/http"
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	urlmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/eventstore/inmemory"
)

var _ = Describe("Application / HTTP", func() {
	var (
		ctrl                       *gomock.Controller
		metrics                    *urlmocks.MockMetrics
		r                          *testingRouter
		shortURLRepository         event.Repository
		loadBalancerURLsRepository event.Repository
		ctx                        context.Context
	)
	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		log.Default().SetOutput(GinkgoWriter)
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		metrics = urlmocks.NewMockMetrics(ctrl)

		shortURLRepository = event.NewRepository(&url.ShortURL{}, inmemory.NewEventStore(), event.NewBroker())
		loadBalancerURLsRepository = event.NewRepository(&url.LoadBalancedURL{}, inmemory.NewEventStore(), event.NewBroker())
		r = newTestingRouter(http.Config{
			BaseDomain:                 "http://example.com",
			ShortURLRepository:         shortURLRepository,
			LoadBalancedURLsRepository: loadBalancerURLsRepository,
			CustomMetrics:              metrics,
		})

		metrics.EXPECT().RecordFileURLMetrics().AnyTimes()
		metrics.EXPECT().RecordSingleURLMetrics().AnyTimes()
	})
	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when it receives an HTTP request for a short URL", func() {
		It("returns the short URL", func() {
			response := r.doPOSTRequest("/api/v1/link", longURLRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(response).To(HaveHTTPBody(MatchJSON(longURLResponse())))

			entity, _, err := shortURLRepository.Load(ctx, "lxqrJ9xF")
			Expect(err).ToNot(HaveOccurred())

			shortURL, ok := entity.(*url.ShortURL)
			Expect(ok).To(BeTrue())
			Expect(shortURL.OriginalURL.URL).To(Equal("https://google.es"))
		})

		Context("but the JSON is malformed", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/v1/link", badURLRequestWithMalformedJSON())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
				Expect(response).To(HaveHTTPBody(ContainSubstring("empty URL requested")))
			})
		})
	})

	Context("when it receives an HTTP request for a load-balancing URL creation", func() {
		It("returns the load balanced URL", func() {
			response := r.doPOSTRequest("/api/v1/loadbalancer", loadBalancerURLRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(response).To(HaveHTTPBody(MatchJSON(loadBalancerURLResponse())))

			entity, _, err := loadBalancerURLsRepository.Load(ctx, "5XEOqhb0")
			Expect(err).ToNot(HaveOccurred())

			loadBalancedURL, ok := entity.(*url.LoadBalancedURL)
			Expect(ok).To(BeTrue())
			Expect(loadBalancedURL.LongURLs).To(ConsistOf(
				url.OriginalURL{URL: "https://google.es", IsValid: false},
				url.OriginalURL{URL: "https://youtube.com", IsValid: false},
			))
		})

		Context("but the JSON is malformed", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/v1/loadbalancer", badURLRequestWithMalformedJSON())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
				Expect(response).To(HaveHTTPBody(ContainSubstring("no URLs specified")))
			})
		})

		Context("but the list of URLs is empty", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/v1/loadbalancer", badLoadBalancerEmptyURLList())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
				Expect(response).To(HaveHTTPBody(ContainSubstring("no URLs specified")))
			})
		})
	})

	Context("when it receives an HTTP request for a redirection", func() {
		Context("and the URL is present in the repository", func() {
			It("responds with a URL redirect", func() {
				// create short url
				r.doPOSTRequest("/api/v1/link", longURLRequest())
				// and verify it TODO(fede): This should be verified through domain event from the message broker
				err := shortURLRepository.Save(ctx, &url.ShortURLVerified{
					Base: event.Base{
						ID:      "lxqrJ9xF",
						Version: 1,
						At:      time.Now(),
					},
				})
				Expect(err).ToNot(HaveOccurred())

				// a redirect is requested
				response := r.doGETRequest("/r/lxqrJ9xF")

				// and the redirection is performed
				Expect(response.StatusCode).To(Equal(gohttp.StatusPermanentRedirect))
				Expect(response.Header.Get("Location")).To(Equal("https://google.es"))

				// the entity saved in the Event Sourcing repository should have 1 click, with entity version 2
				entity, version, err := shortURLRepository.Load(ctx, "lxqrJ9xF")
				Expect(err).ToNot(HaveOccurred())
				Expect(version).To(Equal(2))

				shortURL, ok := entity.(*url.ShortURL)
				Expect(ok).To(BeTrue())
				Expect(shortURL.Clicks).To(Equal(1))
			})
		})

		Context("but the URL is not present in the repository", func() {
			It("returns a 404 error", func() {
				response := r.doGETRequest("/r/123456")

				Expect(response.StatusCode).To(Equal(gohttp.StatusNotFound))
			})
		})
	})

	Context("when it receives an HTTP request for a load-balancing redirection", func() {
		Context("and the URL is present in the repository", func() {
			It("responds with a URL redirect", func() {
				r.doPOSTRequest("/api/v1/loadbalancer", loadBalancerURLRequest())

				err := loadBalancerURLsRepository.Save(ctx, &url.LoadBalancedURLVerified{
					Base: event.Base{
						ID:      "5XEOqhb0",
						Version: 1,
						At:      time.Now(),
					},
					VerifiedURL: "https://google.es",
				})
				Expect(err).ToNot(HaveOccurred())

				response := r.doGETRequest("/lb/5XEOqhb0")

				Expect(response).To(HaveHTTPStatus(gohttp.StatusTemporaryRedirect))
				Expect(response).To(HaveHTTPHeaderWithValue("Location", "https://google.es"))
			})
		})
		Context("and there are multiple valid URLs", func() {
			It("responds with a different URL each time", func() {
				r.doPOSTRequest("/api/v1/loadbalancer", loadBalancerURLRequest())

				err := loadBalancerURLsRepository.Save(ctx,
					&url.LoadBalancedURLVerified{
						Base: event.Base{
							ID:      "5XEOqhb0",
							Version: 1,
							At:      time.Now(),
						},
						VerifiedURL: "https://google.es",
					},
					&url.LoadBalancedURLVerified{
						Base: event.Base{
							ID:      "5XEOqhb0",
							Version: 1,
							At:      time.Now(),
						},
						VerifiedURL: "https://youtube.com",
					},
				)
				Expect(err).ToNot(HaveOccurred())

				Eventually(func() *gohttp.Response {
					return r.doGETRequest("/lb/5XEOqhb0")
				}).Should(HaveHTTPHeaderWithValue("Location", "https://google.es"))
				Eventually(func() *gohttp.Response {
					return r.doGETRequest("/lb/5XEOqhb0")
				}).Should(HaveHTTPHeaderWithValue("Location", "https://youtube.com"))
			})
		})

		Context("but the URL does not have any original valid URL", func() {
			It("returns a 404 error", func() {
				r.doPOSTRequest("/api/v1/loadbalancer", loadBalancerURLRequest())

				response := r.doGETRequest("/lb/5XEOqhb0")

				Expect(response).To(HaveHTTPStatus(gohttp.StatusNotFound))
			})
		})
		Context("but the URL is not present in the repository", func() {
			It("returns a 404 error", func() {
				response := r.doGETRequest("/lb/123456")

				Expect(response).To(HaveHTTPStatus(gohttp.StatusNotFound))
			})
		})
	})

	Context("when it receives an HTTP request to shorten a CSV file", func() {
		It("returns a CSV with the URLs shortened", func() {
			response := r.doPOSTFormRequest("/csv", csvFileRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusCreated))
			Expect(response.Header.Get("Content-type")).To(Equal("text/csv"))
			Expect(response.Header.Get("Location")).To(Equal("google.com"))
			Expect(response).To(HaveHTTPBody(Equal(csvFileResponse())))

			entity, version, err := shortURLRepository.Load(ctx, "uuqVS5Vz")
			Expect(err).ToNot(HaveOccurred())
			Expect(version).To(Equal(0))
			firstURL, ok := entity.(*url.ShortURL)
			Expect(ok).To(BeTrue())

			entity, version, err = shortURLRepository.Load(ctx, "1+IiyNe6")
			Expect(err).ToNot(HaveOccurred())
			Expect(version).To(Equal(0))
			secondURL, ok := entity.(*url.ShortURL)
			Expect(ok).To(BeTrue())

			Expect(firstURL.OriginalURL.URL).To(Equal("google.com"))
			Expect(secondURL.OriginalURL.URL).To(Equal("youtube.com"))
		})

		Context("but the CSV is empty", func() {
			It("returns a bad request code", func() {
				response := r.doPOSTFormRequest("/csv", bytes.NewReader([]byte("")))

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})

		Context("but form-data name field is not equal to 'file'", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTFormRequest("/csv", badCsvFileRequest())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
				Expect(response).To(HaveHTTPBody(ContainSubstring("unable to convert data to long urls: the list of URLs is empty")))
			})
		})
	})

	Context("when it retrieves an HTTP request to an unknown endpoint", func() {
		It("returns an 404 error", func() {
			response := r.doGETRequest("/unknown/endpoint")

			Expect(response.StatusCode).To(Equal(gohttp.StatusNotFound))
		})
	})
})

func csvFileRequest() io.Reader {
	return strings.NewReader(`--unaCadenaDelimitadora
Content-Disposition: form-data; name="file"
Content-Type: text/csv

google.com
youtube.com
--unaCadenaDelimitadora--`)
}

func badCsvFileRequest() io.Reader {
	return strings.NewReader(`--unaCadenaDelimitadora
Content-Disposition: form-data; name="badName"
Content-Type: text/csv

google.com
youtube.com
--unaCadenaDelimitadora--`)
}

func csvFileResponse() []byte {
	return []byte(`google.com,http://example.com/r/uuqVS5Vz,
youtube.com,http://example.com/r/1+IiyNe6,
`)
}

func longURLRequest() io.Reader {
	return strings.NewReader(`{
	"url": "https://google.es"
}`)
}

func loadBalancerURLRequest() io.Reader {
	return strings.NewReader(`{
	"urls": ["https://google.es", "https://youtube.com"]
}`)
}

func loadBalancerURLResponse() string {
	return `{
	"url": "http://example.com/lb/5XEOqhb0"
}`
}

func badLoadBalancerEmptyURLList() io.Reader {
	return strings.NewReader(`{"urls": []}`)
}

func longURLResponse() string {
	return `{
	"url": "http://example.com/r/lxqrJ9xF"
}`
}

func badURLRequestWithMalformedJSON() io.Reader {
	return strings.NewReader(`
{
	"badjson": "https://google.es"
}`)
}
