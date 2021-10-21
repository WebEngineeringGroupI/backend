package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Engine struct {
	urlShortenerService *url.Shortener
	httpDomain          string
}

func (e *Engine) Shortener() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var dataIn shortURLDataIn
		err := json.NewDecoder(request.Body).Decode(&dataIn)
		if err != nil || dataIn.URL == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := e.urlShortenerService.HashFromURL(dataIn.URL)
		if shortURL == nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		dataOut := shortURLDataOut{
			URL: fmt.Sprintf("%s/%s", e.httpDomain, shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func NewEngine(httpDomain string, shortURLRepository url.ShortURLRepository) *Engine {
	return &Engine{
		httpDomain:          httpDomain,
		urlShortenerService: url.NewShortener(shortURLRepository),
	}
}
