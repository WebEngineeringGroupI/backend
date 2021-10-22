package http_test

import (
	"io"
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
		r                  *testingRouter
		shortURLRepository url.ShortURLRepository
	)
	BeforeEach(func() {
		shortURLRepository = inmemory.NewRepository()
		r = newTestingRouter(http.Config{
			BaseDomain:         "http://example.com",
			ShortURLRepository: shortURLRepository,
		})
	})

	Context("when it retrieves an HTTP request for a short URL", func() {
		It("returns the short URL", func() {
			response := r.doPOSTRequest("/api/link", longURLRequest())

			Expect(response.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(readAll(response.Body)).To(MatchJSON(longURLResponse()))

			shortURL, err := shortURLRepository.FindByHash("lxqrJ9xF")

			Expect(err).To(Succeed())
			Expect(shortURL.LongURL).To(Equal("https://google.es"))
		})

		Context("but the JSON is malformed", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/link", badURLRequestWithMalformedJSON())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})
		Context("but the long URL is malformed and not supported", func() {
			It("returns StatusBadRequest code", func() {
				response := r.doPOSTRequest("/api/link", badURLRequestWithFTP())

				Expect(response.StatusCode).To(Equal(gohttp.StatusBadRequest))
			})
		})
	})

	Context("when it retrieves an HTTP request for a redirection", func() {
		Context("and the URL is present in the repository", func() {
			It("responds with a URL redirect", func() {
				_ = shortURLRepository.Save(&url.ShortURL{
					Hash:    "123456",
					LongURL: "https://google.com",
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
	Context("when it retrieves an HTTP request to an unknown endpoint", func() {
		It("returns an 404 error", func() {
			response := r.doGETRequest("/unknown/endpoint")

			Expect(response.StatusCode).To(Equal(gohttp.StatusNotFound))
		})
	})
})

func longURLRequest() io.Reader {
	return strings.NewReader(`{
	"url": "https://google.es"
}`)
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
