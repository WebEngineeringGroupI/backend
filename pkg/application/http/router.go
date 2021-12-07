package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
)

type Config struct {
	BaseDomain                 string
	ShortURLRepository         url.ShortURLRepository
	URLValidator               url.Validator
	CustomMetrics              url.Metrics
	LoadBalancedURLsRepository url.LoadBalancedURLsRepository
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
	h := NewHandlerRepository(config, httprouterVariableExtractor())

	router.Handler(http.MethodPost, "/api/v1/link", h.shortener())
	router.Handler(http.MethodPost, "/api/v1/loadbalancer", h.loadBalancingURLCreator())
	router.Handler(http.MethodPost, "/csv", h.csvShortener())
	router.Handler(http.MethodGet, "/r/:hash", h.redirector())
	router.Handler(http.MethodGet, "/lb/:hash", h.loadBalancingRedirector())
	router.Handler(http.MethodGet, "/metrics", promhttp.Handler())

	router.NotFound = h.notFound()
}
