package http_test

import (
	//"github.com/gorilla/websocket"
	"io"
	gohttp "net/http"
	"net/http/httptest"

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

func (t *testingRouter) doWebSocketRequest(path string) *gohttp.Response {
	//TODO(fede): Set fake ws client
	//_ ,response, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/link",nil)
	//ExpectWithOffset(1, err).To(Succeed())
	//return response
	return nil
}


func newTestingRouter(config http.Config) *testingRouter {
	return &testingRouter{
		config: config,
	}
}
