package http_test

import (
	"bytes"
	"errors"
	"io"
	"log"
	gohttp "net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Application / HTTP", func() {
	var (
		r                          *testingRouter
		shortURLRepository         url.ShortURLRepository
		loadBalancerURLsRepository url.LoadBalancedURLsRepository
		validator                  *FakeURLValidator
		metrics                    *FakeMetrics
	)
	BeforeEach(func() {
		log.Default().SetOutput(GinkgoWriter)
		inmemoryRepository := inmemory.NewRepository()
		shortURLRepository = inmemoryRepository
		loadBalancerURLsRepository = inmemoryRepository
		validator = &FakeURLValidator{returnValidURL: true}
		metrics = &FakeMetrics{}
		r = newTestingRouter(http.Config{
			BaseDomain:                 "http://example.com",
			ShortURLRepository:         shortURLRepository,
			LoadBalancedURLsRepository: loadBalancerURLsRepository,
			URLValidator:               validator,
			CustomMetrics:              metrics,
		})
	})

	Context("when it receives an HTTP request for a short URL", func() {
		It("returns the short URL", func() {
			response := r.doPOSTRequest("/api/v1/link", longURLRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(readAll(response.Body)).To(MatchJSON(longURLResponse()))

			shortURL, err := shortURLRepository.FindShortURLByHash("lxqrJ9xF")

			Expect(err).To(Succeed())
			Expect(shortURL.OriginalURL.URL).To(Equal("https://google.es"))
		})

		Context("but the JSON is malformed", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/v1/link", badURLRequestWithMalformedJSON())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})
		Context("but the long URL is invalid", func() {
			It("returns StatusBadRequest code", func() {
				validator.shouldReturnValidURL(false)
				response := r.doPOSTRequest("/api/v1/link", badURLRequestWithFTP())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})
		Context("but the validator is unable to validate the URL", func() {
			It("returns InternalServerError", func() {
				validator.shouldReturnError(errors.New("error validating the URL"))
				response := r.doPOSTRequest("/api/v1/link", badURLRequestWithFTP())

				Expect(response.StatusCode).To(Equal(gohttp.StatusInternalServerError))
			})
		})
	})

	Context("when it receives an HTTP request for a load-balancing URL creation", func() {
		It("returns the load balanced URL", func() {
			response := r.doPOSTRequest("/api/v1/loadbalancer", loadBalancerURLRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(readAll(response.Body)).To(MatchJSON(loadBalancerURLResponse()))

			loadBalancedURL, err := loadBalancerURLsRepository.FindLoadBalancedURLByHash("5XEOqhb0")
			Expect(err).ToNot(HaveOccurred())
			Expect(loadBalancedURL.LongURLs).To(ConsistOf(
				url.OriginalURL{URL: "https://google.es"},
				url.OriginalURL{URL: "https://youtube.com"},
			))
		})
	})

	Context("when it receives an HTTP request for a redirection", func() {
		Context("and the URL is present in the repository", func() {
			It("responds with a URL redirect", func() {
				_ = shortURLRepository.SaveShortURL(&url.ShortURL{
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

	Context("when it receives an HTTP request to shorten a CSV file", func() {
		It("returns a CSV with the URLs shortened", func() {
			response := r.doPOSTFormRequest("/csv", csvFileRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusCreated))
			Expect(response.Header.Get("Content-type")).To(Equal("text/csv"))
			Expect(response.Header.Get("Location")).To(Equal("google.com"))
			Expect(readAll(response.Body)).To(Equal(csvFileResponse()))

			firstURL, err := shortURLRepository.FindShortURLByHash("uuqVS5Vz")
			Expect(err).To(Succeed())
			secondURL, err := shortURLRepository.FindShortURLByHash("1+IiyNe6")
			Expect(err).To(Succeed())

			Expect(firstURL.OriginalURL.URL).To(Equal("google.com"))
			Expect(secondURL.OriginalURL.URL).To(Equal("youtube.com"))
		})

		Context("but the CSV is empty", func() {
			It("returns a bad request code", func() {
				response := r.doPOSTFormRequest("/csv", bytes.NewReader([]byte("")))

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})

		Context("but a long URL is invalid", func() {
			It("returns StatusBadRequest code", func() {
				validator.shouldReturnValidURL(false)
				response := r.doPOSTFormRequest("/csv", csvFileRequest())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})
		Context("but form-data name field is not equal to 'file'", func() {
			It("returns StatusBadRequest code", func() {
				validator.shouldReturnValidURL(true)
				response := r.doPOSTFormRequest("/csv", badCsvFileRequest())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
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

func csvFileResponse() string {
	return `google.com,http://example.com/r/uuqVS5Vz,
youtube.com,http://example.com/r/1+IiyNe6,
`
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

func badURLRequestWithFTP() io.Reader {
	return strings.NewReader(`
{
	"url": "ftp://google.es"
}`)
}

func readAll(reader io.Reader) string {
	bytes, err := io.ReadAll(reader)

	ExpectWithOffset(1, err).To(Succeed())
	return string(bytes)
}
