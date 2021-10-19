package http_test

import (
	`io`
	gohttp `net/http`
	`net/http/httptest`
	`strings`

	. "github.com/onsi/ginkgo"
	. `github.com/onsi/gomega`

	`github.com/WebEngineeringGroupI/backend/pkg/application/http`
	`github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory`
)

var _ = Describe("Application / HTTP", func() {
	Context("when it retrieves an HTTP request for a short URL", func() {
		It("returns the short URL", func() {
			httpEngine := http.NewEngine("http://example.com", inmemory.NewRepository())
			handler := httpEngine.Shortener()

			request := httptest.NewRequest("POST", "http://example.com/api/link", longURLRequest())
			recorder := httptest.NewRecorder()
			handler(recorder, request)
			result := recorder.Result()

			Expect(result.StatusCode).To(Equal(gohttp.StatusOK))
			Expect(readAll(result.Body)).To(MatchJSON(longURLResponse()))
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

func readAll(reader io.Reader) string {
	bytes, err := io.ReadAll(reader)

	ExpectWithOffset(1, err).To(Succeed())
	return string(bytes)
}
