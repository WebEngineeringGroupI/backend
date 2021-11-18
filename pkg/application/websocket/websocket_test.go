package websocket_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/WebEngineeringGroupI/backend/pkg/application/websocket"
	"github.com/WebEngineeringGroupI/backend/pkg/domain"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/inmemory"
)

var _ = Describe("Websocket", func() {
	var (
		handler            websocket.MessageHandler
		shortURLRepository url.ShortURLRepository
		validator          *FakeURLValidator
	)
	BeforeEach(func() {
		wholeURL := domain.NewWholeURL("https://example.com")
		shortURLRepository = inmemory.NewRepository()
		validator = &FakeURLValidator{returnValidURL: true}
		shortener := url.NewSingleURLShortener(shortURLRepository, validator)
		handler = websocket.NewMessageHandler(wholeURL, shortener)
	})

	When("it receives a ping message", func() {
		It("returns the pong message", func() {
			responseType, response := handler.HandleMessage(websocket.PingMessage, []byte("ping"))

			Expect(responseType).To(Equal(websocket.PongMessage))
			Expect(response).To(Equal([]byte("pong")))
		})
	})

	When("it receives a short URL request", func() {
		It("returns the URLs shortened", func() {
			responseType, response := handler.HandleMessage(websocket.TextMessage, urlsToShort())

			Expect(responseType).To(Equal(websocket.TextMessage))
			Expect(response).To(MatchJSON(urlsShorted()))
		})
	})

	When("the validator is unable to validate the URLs", func() {
		It("returns the error in the URLs shortened", func() {
			validator.shouldReturnError(errors.New("unknown error"))
			responseType, response := handler.HandleMessage(websocket.TextMessage, urlsToShort())

			Expect(responseType).To(Equal(websocket.TextMessage))
			Expect(response).To(MatchJSON(urlsShortedWithValidationErrors()))
		})
	})

	When("the validator does not validate the URLs", func() {
		It("returns the invalid URL message in the text", func() {
			validator.shouldReturnValidURL(false)
			responseType, response := handler.HandleMessage(websocket.TextMessage, urlsToShort())

			Expect(responseType).To(Equal(websocket.TextMessage))
			Expect(response).To(MatchJSON(urlsShortedWithInvalidURLs()))
		})
	})
})

func urlsToShort() []byte {
	return []byte(`{
  "request_type": "short_urls",
  "request": {
    "urls": ["https://google.com", "https://youtube.com"]
  }
}`)
}

func urlsShorted() string {
	return `{
  "response_type": "short_urls",
  "response": {
    "urls": ["https://example.com/r/cv6VxVdu", "https://example.com/r/unW6a4Dd"]
  }
}`
}

func urlsShortedWithInvalidURLs() interface{} {
	return `{
  "response_type": "short_urls",
  "response": {
    "urls": ["unable to short URL: invalid long URL specified",
            "unable to short URL: invalid long URL specified"]
  }
}`
}

func urlsShortedWithValidationErrors() interface{} {
	return `{
  "response_type": "short_urls",
  "response": {
    "urls": ["unable to short URL: unknown error", "unable to short URL: unknown error"]
  }
}`
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
