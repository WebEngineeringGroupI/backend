package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Config struct {
	BaseDomain         string
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
	h := NewHandlerRepository(config.BaseDomain, httprouterVariableExtractor())

	router.Handler(http.MethodPost, "/api/v1/link", h.shortener(config.ShortURLRepository, config.URLValidator))
	router.Handler(http.MethodGet, "/ws/link", h.wsshortener(config.ShortURLRepository, config.URLValidator))
	router.Handler(http.MethodPost, "/csv", h.csvShortener(config.ShortURLRepository, config.URLValidator))
	router.Handler(http.MethodGet, "/r/:hash", h.redirector(config.ShortURLRepository, config.URLValidator))

	router.NotFound = h.notFound()
}
