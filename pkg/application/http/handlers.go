package http

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	gorillaws "github.com/gorilla/websocket"

	"github.com/WebEngineeringGroupI/backend/pkg/application/websocket"
	"github.com/WebEngineeringGroupI/backend/pkg/domain"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/redirect"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url/formatter"
)

type HandlerRepository struct {
	variableExtractor VariableExtractor
	wholeURL          *domain.WholeURL
}

type VariableExtractor interface {
	Extract(request *http.Request, key string) string
}

func (e *HandlerRepository) Shortener(repository url.ShortURLRepository, validator url.Validator) http.HandlerFunc {
	urlShortener := url.NewSingleURLShortener(repository, validator)

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
			URL: e.wholeURL.FromHash(shortURL.Hash),
		}
		err = json.NewEncoder(writer).Encode(&dataOut)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("error marshaling the response: %s", err)
			return
		}
	}
}

func (e *HandlerRepository) Redirector(repository url.ShortURLRepository, validator url.Validator) http.HandlerFunc {
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

func (e *HandlerRepository) NotFound() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		http.NotFound(writer, request)
	}
}

func (e *HandlerRepository) CSVShortener(repository url.ShortURLRepository, validator url.Validator) http.HandlerFunc {
	csvShortener := url.NewFileURLShortener(repository, validator, formatter.NewCSV())

	return func(writer http.ResponseWriter, request *http.Request) {
		data, err := io.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "unable to read all request body", http.StatusInternalServerError)
			return
		}

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
				e.wholeURL.FromHash(shortURL.Hash),
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

func (e *HandlerRepository) WSHandler(config Config) http.HandlerFunc {
	urlShortener := url.NewSingleURLShortener(config.ShortURLRepository, config.URLValidator)
	websocketMessageHandler := websocket.NewMessageHandler(config.WholeURL, urlShortener)

	var upgrader = gorillaws.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	return func(writer http.ResponseWriter, request *http.Request) {
		ws, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			http.Error(writer, "error upgrading to websocket connection", http.StatusInternalServerError)
			log.Printf("error upgrading to websocket connection: %s", err)
			return
		}
		defer ws.Close()

		for {
			messageType, message, err := ws.ReadMessage()
			if err != nil {
				log.Printf("error reading message from websocket: %s", err)
				break
			}
			responseType, response := websocketMessageHandler.HandleMessage(messageType, message)
			err = ws.WriteMessage(responseType, response)
			if err != nil {
				log.Printf("error writing message to websocket: %s", err)
				break
			}
		}
	}
}

func NewHandlerRepository(wholeURL *domain.WholeURL, variableExtractor VariableExtractor) *HandlerRepository {
	return &HandlerRepository{
		wholeURL:          wholeURL,
		variableExtractor: variableExtractor,
	}
}
