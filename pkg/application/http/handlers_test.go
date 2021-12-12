package http_test

import (
	"bytes"
	"context"
	"io"
	"log"
	"math/rand"
	gohttp "net/http"
	"strings"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	urlmocks "github.com/WebEngineeringGroupI/backend/pkg/domain/url/mocks"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Application / HTTP", func() {
	var (
		ctrl                       *gomock.Controller
		emitter                    *mocks.MockEmitter
		metrics                    *urlmocks.MockMetrics
		r                          *testingRouter
		shortURLRepository         url.ShortURLRepository
		loadBalancerURLsRepository url.LoadBalancedURLsRepository
		ctx                        context.Context
	)
	BeforeEach(func() {
		rand.Seed(GinkgoRandomSeed())
		log.Default().SetOutput(GinkgoWriter)
		ctx = context.Background()

		ctrl = gomock.NewController(GinkgoT())
		metrics = urlmocks.NewMockMetrics(ctrl)
		emitter = mocks.NewMockEmitter(ctrl)

		inmemoryRepository := inmemory.NewRepository()
		shortURLRepository = inmemoryRepository
		loadBalancerURLsRepository = inmemoryRepository
		r = newTestingRouter(http.Config{
			BaseDomain:                 "http://example.com",
			ShortURLRepository:         shortURLRepository,
			LoadBalancedURLsRepository: loadBalancerURLsRepository,
			CustomMetrics:              metrics,
			EventEmitter:               emitter,
		})

		emitter.EXPECT().EmitShortURLCreated(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
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

			shortURL, err := shortURLRepository.FindShortURLByHash(ctx, "lxqrJ9xF")

			Expect(err).ToNot(HaveOccurred())
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

			loadBalancedURL, err := loadBalancerURLsRepository.FindLoadBalancedURLByHash(ctx, "5XEOqhb0")
			Expect(err).ToNot(HaveOccurred())
			Expect(loadBalancedURL.LongURLs).To(ConsistOf(
				url.OriginalURL{URL: "https://google.es"},
				url.OriginalURL{URL: "https://youtube.com"},
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
				_ = shortURLRepository.SaveShortURL(ctx, &url.ShortURL{
					Hash:        "123456",
					OriginalURL: url.OriginalURL{URL: "https://google.com", IsValid: true},
				})

				response := r.doGETRequest("/r/123456")

				Expect(response.StatusCode).To(Equal(gohttp.StatusPermanentRedirect))
				Expect(response.Header.Get("Location")).To(Equal("https://google.com"))
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
				_ = loadBalancerURLsRepository.SaveLoadBalancedURL(ctx, &url.LoadBalancedURL{
					Hash: "123456",
					LongURLs: []url.OriginalURL{
						{URL: "https://google.com", IsValid: true},
						{URL: "https://youtube.com", IsValid: false},
					},
				})

				response := r.doGETRequest("/lb/123456")

				Expect(response).To(HaveHTTPStatus(gohttp.StatusTemporaryRedirect))
				Expect(response).To(HaveHTTPHeaderWithValue("Location", "https://google.com"))
			})
		})
		Context("and there are multiple valid URLs", func() {
			It("responds with a different URL each time", func() {
				_ = loadBalancerURLsRepository.SaveLoadBalancedURL(ctx, &url.LoadBalancedURL{
					Hash: "123456",
					LongURLs: []url.OriginalURL{
						{URL: "https://google.com", IsValid: true},
						{URL: "https://youtube.com", IsValid: true},
					},
				})

				Eventually(func() *gohttp.Response {
					return r.doGETRequest("/lb/123456")
				}).Should(HaveHTTPHeaderWithValue("Location", "https://google.com"))
				Eventually(func() *gohttp.Response {
					return r.doGETRequest("/lb/123456")
				}).Should(HaveHTTPHeaderWithValue("Location", "https://youtube.com"))
			})
		})

		Context("but the URL does not have any original valid URL", func() {
			It("returns a 404 error", func() {
				_ = loadBalancerURLsRepository.SaveLoadBalancedURL(ctx, &url.LoadBalancedURL{
					Hash: "123456",
					LongURLs: []url.OriginalURL{
						{URL: "https://google.com", IsValid: false},
						{URL: "https://youtube.com", IsValid: false},
					},
				})

				response := r.doGETRequest("/lb/123456")

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

			firstURL, err := shortURLRepository.FindShortURLByHash(ctx, "uuqVS5Vz")
			Expect(err).ToNot(HaveOccurred())
			secondURL, err := shortURLRepository.FindShortURLByHash(ctx, "1+IiyNe6")
			Expect(err).ToNot(HaveOccurred())

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
