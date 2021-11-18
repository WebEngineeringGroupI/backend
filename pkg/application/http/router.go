package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/WebEngineeringGroupI/backend/pkg/domain"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Config struct {
	WholeURL           *domain.WholeURL
	ShortURLRepository url.ShortURLRepository
	URLValidator       url.Validator
}

func NewRouter(config Config) http.Handler {
	router := httprouter.New()
	registerPaths(router, config)

	return cors.Default().Handler(router)
}

type variableExtractorFunc func(request *http.Request, key string) string

func (v variableExtractorFunc) Extract(request *http.Request, key string) string {
	return v(request, key)
}

func httprouterVariableExtractor() variableExtractorFunc {
	return func(request *http.Request, key string) string {
		return httprouter.ParamsFromContext(request.Context()).ByName(key)
	}
}

func registerPaths(router *httprouter.Router, config Config) {
	restHandler := NewHandlerRepository(config.WholeURL, httprouterVariableExtractor())

	router.Handler(http.MethodPost, "/api/v1/link", restHandler.Shortener(config.ShortURLRepository, config.URLValidator))
	router.Handler(http.MethodPost, "/csv", restHandler.CSVShortener(config.ShortURLRepository, config.URLValidator))
	router.Handler(http.MethodGet, "/r/:hash", restHandler.Redirector(config.ShortURLRepository, config.URLValidator))

	router.Handler(http.MethodGet, "/ws/link", restHandler.WSHandler(config))

	router.NotFound = restHandler.NotFound()
}
