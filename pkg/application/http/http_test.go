package http_test

import (
	"io"
	gohttp "net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Application / HTTP", func() {
	var (
		httpEngine *http.Engine
	)

	BeforeEach(func() {
		httpEngine = http.NewEngine("http://example.com", inmemory.NewRepository())
	})
	Context("when it retrieves an HTTP request for a short URL", func() {
		It("returns the short URL", func() {
			handler := httpEngine.Shortener()
			request := httptest.NewRequest("POST", "http://example.com/api/link", longURLRequest())
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()

			Expect(result.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(readAll(result.Body)).To(MatchJSON(longURLResponse()))
		})
	})
	Context("when it retrieves an HTTP request for a short URL with malformed JSON key ", func() {
		It("returns StatusBadRequest code", func() {
			handler := httpEngine.Shortener()
			request := httptest.NewRequest("POST", "http://example.com/api/link", badjsonURLRequest())
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()

			Expect(result.StatusCode).To(Equal(gohttp.StatusBadRequest))

		})
	})
	Context("when it retrieves an HTTP request for a short URL with malformed long URL", func() {
		It("returns StatusBadRequest code", func() {
			handler := httpEngine.Shortener()
			request := httptest.NewRequest("POST", "http://example.com/api/link", badURLRequest())
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()

			Expect(result.StatusCode).To(Equal(gohttp.StatusBadRequest))
		})
	})
})

func longURLRequest() io.Reader {
	return strings.NewReader(`
{
	"url": "https://google.es"
}
`)
}

func longURLResponse() string {
	return `
{
	"url": "http://example.com/lxqrJ9xF"
}
`
}

func badjsonURLRequest() io.Reader {
	return strings.NewReader(`
{
	"badjson": "https://google.es"
}
`)
}

func badURLRequest() io.Reader {
	return strings.NewReader(`
{
	"url": "ftp://google.es"
}
`)
}

func readAll(reader io.Reader) string {
	bytes, err := io.ReadAll(reader)

	ExpectWithOffset(1, err).To(Succeed())
	return string(bytes)
}
