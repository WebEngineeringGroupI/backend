package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type HandlerRepository struct {
	baseDomain        string
	variableExtractor VariableExtractor
}

type VariableExtractor interface {
	Extract(request *http.Request, key string) string
}

func (e *HandlerRepository) shortener(repository url.ShortURLRepository) http.HandlerFunc {
	urlShortener := url.NewShortener(repository)

	return func(writer http.ResponseWriter, request *http.Request) {
		var dataIn shortURLDataIn
		err := json.NewDecoder(request.Body).Decode(&dataIn)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if dataIn.URL == "" {
			http.Error(writer, "empty URL requested", http.StatusBadRequest)
			return
		}

		shortURL, err := urlShortener.HashFromURL(dataIn.URL)
		if errors.Is(err, url.ErrInvalidLongURLSpecified) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			log.Printf("error retrieving hash from long URL: %s", err)
			return
		}

		dataOut := shortURLDataOut{
			URL: fmt.Sprintf("%s/r/%s", e.baseDomain, shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshaling the response: %s", err)
			return
		}
	}
}

func (e *HandlerRepository) redirector(repository url.ShortURLRepository) http.HandlerFunc {
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

func (e *HandlerRepository) notFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	}
}

func NewHandlerRepository(baseDomain string, variableExtractor VariableExtractor) *HandlerRepository {
	return &HandlerRepository{
		baseDomain:        strings.TrimSuffix(baseDomain, "/"),
		variableExtractor: variableExtractor,
	}
}
