package http_test

import (
	"io"
	gohttp "net/http"
	"net/http/httptest"

	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
)

type testingRouter struct {
	config http.Config
}

func (t *testingRouter) doPOSTFormRequest(path string, body io.Reader) *gohttp.Response {
	request, err := gohttp.NewRequest(gohttp.MethodPost, path, body)
	request.Header.Set("Content-Type", "multipart/form-data; boundary=unaCadenaDelimitadora")
	ExpectWithOffset(1, err).To(Succeed())

	recorder := httptest.NewRecorder()
	router := http.NewRouter(t.config)
	router.ServeHTTP(recorder, request)

	return recorder.Result()
}

func (t *testingRouter) doPOSTRequest(path string, body io.Reader) *gohttp.Response {
	request, err := gohttp.NewRequest(gohttp.MethodPost, path, body)
	ExpectWithOffset(1, err).To(Succeed())

	recorder := httptest.NewRecorder()
	router := http.NewRouter(t.config)
	router.ServeHTTP(recorder, request)

	return recorder.Result()
}

func (t *testingRouter) doGETRequest(path string) *gohttp.Response {
	request, err := gohttp.NewRequest(gohttp.MethodGet, path, nil)
	ExpectWithOffset(1, err).To(Succeed())

	recorder := httptest.NewRecorder()
	router := http.NewRouter(t.config)
	router.ServeHTTP(recorder, request)

	return recorder.Result()
}

func newTestingRouter(config http.Config) *testingRouter {
	return &testingRouter{
		config: config,
	}
}

type FakeURLValidator struct {
	returnValidURL bool
	returnError    error
}

func (f *FakeURLValidator) shouldReturnValidURL(validURL bool) {
	f.returnValidURL = validURL
}

func (f *FakeURLValidator) shouldReturnError(err error) {
	f.returnError = err
}

func (f *FakeURLValidator) ValidateURL(url string) (bool, error) {
	return f.returnValidURL, f.returnError
}

func (f *FakeURLValidator) ValidateURLs(urls []string) (bool, error) {
	return f.returnValidURL, f.returnError
}

type FakeMetrics struct {
	singleURLMetrics   int
	multipleURLMetrics int
	fileURLMetrics     int
	urlsProcessed      int
}

func (f *FakeMetrics) RecordSingleURLMetrics() {
	f.singleURLMetrics++
}

func (f *FakeMetrics) RecordFileURLMetrics() {
	f.fileURLMetrics++
}

func (f *FakeMetrics) RecordUrlsProcessed() {
	f.urlsProcessed++
}
