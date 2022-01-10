package main

import (
	"log"
	gohttp "net/http"
	"os"
	"strconv"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	gogrpc "google.golang.org/grpc"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres/serializer/json"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/metrics"
)

type factory struct {
	metricsSingleton     url.Metrics
	eventBrokerSingleton event.Broker
}

func (f *factory) NewHTTPAndGRPCWebRouter() gohttp.Handler {
	httpRouter := http.NewRouter(f.httpConfig())
	grpcWebServer := grpcweb.WrapServer(f.NewGRPCServer(),
		grpcweb.WithWebsockets(true),
		grpcweb.WithOriginFunc(func(origin string) bool {
			return true
		}),
	)
	return gohttp.HandlerFunc(func(writer gohttp.ResponseWriter, request *gohttp.Request) {
		if grpcWebServer.IsAcceptableGrpcCorsRequest(request) || grpcWebServer.IsGrpcWebRequest(request) {
			grpcWebServer.ServeHTTP(writer, request)
			return
		}
		httpRouter.ServeHTTP(writer, request)
	})
}

func (f *factory) NewGRPCServer() *gogrpc.Server {
	return grpc.NewServer(f.grpcConfig())
}

func (f *factory) httpConfig() http.Config {
	return http.Config{
		BaseDomain:                 f.baseDomain(),
		CustomMetrics:              f.customMetrics(),
		ShortURLRepository:         f.newShortURLRepository(),
		LoadBalancedURLsRepository: f.newLoadBalancedURLsRepository(),
	}
}

func (f *factory) grpcConfig() grpc.Config {
	return grpc.Config{
		BaseDomain:                 f.baseDomain(),
		ShortURLRepository:         f.newShortURLRepository(),
		CustomMetrics:              f.customMetrics(),
		LoadBalancedURLsRepository: f.newLoadBalancedURLsRepository(),
	}
}

func (f *factory) customMetrics() url.Metrics {
	if f.metricsSingleton == nil {
		f.metricsSingleton = metrics.NewPrometheusMetrics()
	}
	return f.metricsSingleton
}

func (f *factory) baseDomain() string {
	baseDomain, isSet := os.LookupEnv("BASE_DOMAIN")
	if !isSet {
		return "http://localhost:8080"
	}
	return strings.TrimSuffix(baseDomain, "/")
}

func (f *factory) postgresConnectionDetails() *postgres.ConnectionDetails {
	dbPort, err := strconv.Atoi(f.mandatoryEnvVarValue("DB_PORT"))
	if err != nil {
		log.Fatalf("unable to parse DB_PORT, make sure it is defined and is a valid integer")
	}

	return &postgres.ConnectionDetails{
		User:     f.mandatoryEnvVarValue("DB_USER"),
		Pass:     f.mandatoryEnvVarValue("DB_PASS"),
		Host:     f.mandatoryEnvVarValue("DB_HOST"),
		Port:     dbPort,
		Database: f.mandatoryEnvVarValue("DB_NAME"),
		SSLMode:  f.mandatoryEnvVarValue("DB_SSL_MODE"),
	}
}

func (f *factory) mandatoryEnvVarValue(variable string) string {
	value, isSet := os.LookupEnv(variable)
	if !isSet {
		log.Fatalf("mandatory %s env var is not set", variable)
	}
	return value
}

func (f *factory) newPostgresDB(eventSerializer serializer.Serializer) *postgres.DB {
	db, err := postgres.NewDB(f.postgresConnectionDetails(), eventSerializer)
	if err != nil {
		log.Fatalf("unable to create the database connection: %s", err)
	}
	return db
}

func (f *factory) newShortURLRepository() event.Repository {
	return event.NewRepository(&url.ShortURL{}, f.newPostgresDB(json.NewSerializer(
		&url.ShortURLCreated{},
		&url.ShortURLVerified{},
		&url.ShortURLClicked{},
	)), f.eventBroker())
}

func (f *factory) newLoadBalancedURLsRepository() event.Repository {
	return event.NewRepository(&url.LoadBalancedURL{}, f.newPostgresDB(json.NewSerializer(
		&url.LoadBalancedURLCreated{},
		&url.LoadBalancedURLVerified{},
	)), f.eventBroker())
}

func (f *factory) eventBroker() event.Broker {
	if f.eventBrokerSingleton == nil {
		f.eventBrokerSingleton = event.NewBroker()
	}
	return f.eventBrokerSingleton
}

func newFactory() *factory {
	return &factory{}
}
