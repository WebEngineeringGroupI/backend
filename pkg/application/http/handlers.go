package http

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/formatter"
)

type HandlerRepository struct {
	baseDomain        string
	variableExtractor VariableExtractor
}

type VariableExtractor interface {
	Extract(request *http.Request, key string) string
}

func (e *HandlerRepository) shortener(repository url.ShortURLRepository, validator url.Validator, metrics url.Metrics) http.HandlerFunc {
	urlShortener := url.NewSingleURLShortener(repository, validator, metrics)

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
			log.Print(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, url.ErrUnableToValidateURLs) {
			log.Print(err.Error())
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

func (e *HandlerRepository) redirector(repository url.ShortURLRepository, validator url.Validator) http.HandlerFunc {
	redirector := redirect.NewRedirector(repository, validator)

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

func (e *HandlerRepository) csvShortener(repository url.ShortURLRepository, validator url.Validator, metrics url.Metrics) http.HandlerFunc {
	csvShortener := url.NewFileURLShortener(repository, validator, metrics, formatter.NewCSV())

	return func(writer http.ResponseWriter, request *http.Request) {
		data := []byte(request.FormValue("file"))
		shortURLs, err := csvShortener.HashesFromURLData(data)
		if errors.Is(err, url.ErrInvalidLongURLSpecified) || errors.Is(err, url.ErrUnableToConvertDataToLongURLs) {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			log.Printf("error retrieving hash from long URL: %s", err)
			return
		}

		dataOut := csvDataOut{}
		for _, shortURL := range shortURLs {
			dataOut = append(dataOut, []string{
				shortURL.LongURL,
				fmt.Sprintf("%s/r/%s", e.baseDomain, shortURL.Hash),
				"",
			})
		}

		writer.Header().Set("Location", shortURLs[0].LongURL)
		writer.Header().Set("Content-type", "text/csv")
		writer.WriteHeader(http.StatusCreated)
		err = csv.NewWriter(writer).WriteAll(dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshaling the response: %s", err)
			return
		}
	}
}

func NewHandlerRepository(baseDomain string, variableExtractor VariableExtractor) *HandlerRepository {
	return &HandlerRepository{
		baseDomain:        strings.TrimSuffix(baseDomain, "/"),
		variableExtractor: variableExtractor,
	}
}
