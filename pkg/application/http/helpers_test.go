package http_test

import (
	"fmt"
	"io"
	gohttp "net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/websocket"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
)

type testingRouter struct {
	config http.Config
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

func (t *testingRouter) doWebSocketHandshake(path string) (*gohttp.Response, error) {
	server := httptest.NewServer(http.NewRouter(t.config))
	defer server.Close()
	serverURL := fmt.Sprintf("%s/%s", strings.Replace(server.URL, "http", "ws", 1), strings.TrimPrefix(path, "/"))

	conn, response, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return response, nil
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
