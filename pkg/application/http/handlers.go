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
	config            Config
	variableExtractor VariableExtractor
}

type VariableExtractor interface {
	Extract(request *http.Request, key string) string
}

func (e *HandlerRepository) shortener() http.HandlerFunc {
	urlShortener := url.NewSingleURLShortener(e.config.ShortURLRepository, e.config.CustomMetrics, e.config.EventOutbox)

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

		shortURL, err := urlShortener.HashFromURL(request.Context(), dataIn.URL)
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
			URL: fmt.Sprintf("%s/r/%s", e.baseDomain(), shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshaling the response: %s", err)
			return
		}
	}
}

func (e *HandlerRepository) loadBalancingURLCreator() http.HandlerFunc {
	loadBalancerCreator := url.NewLoadBalancer(e.config.LoadBalancedURLsRepository)

	return func(writer http.ResponseWriter, request *http.Request) {
		var dataIn loadBalancerURLDataIn
		err := json.NewDecoder(request.Body).Decode(&dataIn)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		shortURL, err := loadBalancerCreator.ShortURLs(request.Context(), dataIn.URLs)
		if errors.Is(err, url.ErrNoURLsSpecified) {
			log.Print(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if errors.Is(err, url.ErrTooMuchMultipleURLs) {
			log.Print(err.Error())
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		if err != nil {
			http.Error(writer, "internal server error", http.StatusInternalServerError)
			log.Printf("error retrieving hash from long URL: %s", err)
			return
		}

		dataOut := loadBalancerURLDataOut{
			URL: fmt.Sprintf("%s/lb/%s", e.baseDomain(), shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshaling the response: %s", err)
			return
		}
	}
}

func (e *HandlerRepository) redirector() http.HandlerFunc {
	redirector := redirect.NewRedirector(e.config.ShortURLRepository)

	return func(writer http.ResponseWriter, request *http.Request) {
		shortURLHash := e.variableExtractor.Extract(request, "hash")

		originalURL, err := redirector.ReturnOriginalURL(request.Context(), shortURLHash)
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

func (e *HandlerRepository) loadBalancingRedirector() http.HandlerFunc {
	redirector := redirect.NewLoadBalancerRedirector(e.config.LoadBalancedURLsRepository)

	return func(writer http.ResponseWriter, request *http.Request) {
		hash := e.variableExtractor.Extract(request, "hash")

		originalURL, err := redirector.ReturnAValidOriginalURL(request.Context(), hash)
		if errors.Is(err, url.ErrValidURLNotFound) {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(writer, request, originalURL, http.StatusTemporaryRedirect)
	}
}

func (e *HandlerRepository) notFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	}
}

func (e *HandlerRepository) csvShortener() http.HandlerFunc {
	csvShortener := url.NewFileURLShortener(e.config.ShortURLRepository, e.config.CustomMetrics, formatter.NewCSV())

	return func(writer http.ResponseWriter, request *http.Request) {
		data := []byte(request.FormValue("file"))
		shortURLs, err := csvShortener.HashesFromURLData(request.Context(), data)
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
				shortURL.OriginalURL.URL,
				fmt.Sprintf("%s/r/%s", e.baseDomain(), shortURL.Hash),
				"",
			})
		}

		writer.Header().Set("Location", shortURLs[0].OriginalURL.URL)
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

func (e *HandlerRepository) baseDomain() string {
	return strings.TrimSuffix(e.config.BaseDomain, "/")
}

func NewHandlerRepository(config Config, variableExtractor VariableExtractor) *HandlerRepository {
	return &HandlerRepository{
		config:            config,
		variableExtractor: variableExtractor,
	}
}
