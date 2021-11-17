package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/metrics"
)

type Config struct {
	BaseDomain         string
	ShortURLRepository url.ShortURLRepository
	CustomMetrics metrics.CustomMetrics
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

	router.Handler(http.MethodPost, "/api/link", h.shortener(config.ShortURLRepository))
	router.Handler(http.MethodGet, "/r/:hash", h.redirector(config.ShortURLRepository))

	router.NotFound = h.notFound()
}
