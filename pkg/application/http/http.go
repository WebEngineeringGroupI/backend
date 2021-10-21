package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	`strings`

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Engine struct {
	baseDomain        string
	variableExtractor VariableExtractor
}

type VariableExtractor interface {
	Extract(request *http.Request, key string) string
}

func (e *Engine) shortener(repository url.ShortURLRepository) http.HandlerFunc {
	urlShortener := url.NewShortener(repository)

	return func(writer http.ResponseWriter, request *http.Request) {
		var dataIn shortURLDataIn
		err := json.NewDecoder(request.Body).Decode(&dataIn)
		if err != nil || dataIn.URL == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		shortURL := urlShortener.HashFromURL(dataIn.URL)
		if shortURL == nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		dataOut := shortURLDataOut{
			URL: fmt.Sprintf("%s/%s", e.baseDomain, shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		writer.WriteHeader(http.StatusOK)
	}
}

func (e *Engine) redirector(repository url.ShortURLRepository) http.HandlerFunc {
	redirector := redirect.NewRedirector(repository)

	return func(writer http.ResponseWriter, request *http.Request) {
		shortURLHash := e.variableExtractor.Extract(request, "hash")

		originalURL, err := redirector.ReturnOriginalURL(shortURLHash)
		if errors.Is(err, url.ErrShortURLNotFound) {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, originalURL, http.StatusPermanentRedirect)
	}
}

func (e *Engine) notFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	}
}

func NewHandlerRepository(baseDomain string, variableExtractor VariableExtractor) *Engine {
	return &Engine{
		baseDomain:        strings.TrimSuffix(baseDomain, "/"),
		variableExtractor: variableExtractor,
	}
}
