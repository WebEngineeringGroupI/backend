package main

import (
	"log"
	gohttp "net/http"
	"os"
	"strconv"
	"time"

	gogrpc "google.golang.org/grpc"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/metrics"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/pipeline"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/reachable"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/safebrowsing"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/validator/schema"
)

type factory struct {
	shortURLRepositorySingleton         url.ShortURLRepository
	urlValidatorSingleton               url.Validator
	metricsSingleton                    url.Metrics
	loadBalancedURLsRepositorySingleton url.LoadBalancedURLsRepository
	postgresDBSingleton                 *postgres.DB
}

func (f *factory) NewHTTPRouter() gohttp.Handler {
	return http.NewRouter(f.httpConfig())
}

func (f *factory) NewGRPCServer() *gogrpc.Server {
	return grpc.NewServer(f.grpcConfig())
}

func (f *factory) httpConfig() http.Config {
	return http.Config{
		BaseDomain:                 f.baseDomain(),
		ShortURLRepository:         f.shortURLRepository(),
		URLValidator:               f.urlValidator(),
		CustomMetrics:              f.customMetrics(),
		LoadBalancedURLsRepository: f.loadBalancedURLsRepository(),
	}
}

func (f *factory) grpcConfig() grpc.Config {
	return grpc.Config{
		BaseDomain:                 f.baseDomain(),
		ShortURLRepository:         f.shortURLRepository(),
		URLValidator:               f.urlValidator(),
		CustomMetrics:              f.customMetrics(),
		LoadBalancedURLsRepository: f.loadBalancedURLsRepository(),
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
	return baseDomain
}

func (f *factory) shortURLRepository() url.ShortURLRepository {
	if f.shortURLRepositorySingleton == nil {
		f.shortURLRepositorySingleton = f.newPostgresDB()
	}
	return f.shortURLRepositorySingleton
}

func (f *factory) postgresConnectionDetails() postgres.ConnectionDetails {
	dbPort, err := strconv.Atoi(f.mandatoryEnvVarValue("DB_PORT"))
	if err != nil {
		log.Fatalf("unable to parse DB_PORT, make sure it is defined and is a valid integer")
	}

	return postgres.ConnectionDetails{
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

func (f *factory) urlValidator() url.Validator {
	if f.urlValidatorSingleton == nil {
		schemaValidator := schema.NewValidator("http", "https")
		reachableValidator := reachable.NewValidator(gohttp.DefaultClient, 2*time.Second)
		safeBrowsingValidator, err := safebrowsing.NewValidator(f.mandatoryEnvVarValue("SAFE_BROWSING_API_KEY"))
		if err != nil {
			log.Fatalf("unable to build SafeBrowsing URL validator: %s", err)
		}
		f.urlValidatorSingleton = pipeline.NewValidator(schemaValidator, reachableValidator, safeBrowsingValidator)
	}
	return f.urlValidatorSingleton
}

func (f *factory) loadBalancedURLsRepository() url.LoadBalancedURLsRepository {
	if f.loadBalancedURLsRepositorySingleton == nil {
		f.loadBalancedURLsRepositorySingleton = f.newPostgresDB()
	}
	return f.loadBalancedURLsRepositorySingleton
}

func (f *factory) newPostgresDB() *postgres.DB {
	if f.postgresDBSingleton == nil {
		db, err := postgres.NewDB(f.postgresConnectionDetails())
		if err != nil {
			log.Fatalf("unable to create the database connection: %s", err)
		}
		f.postgresDBSingleton = db
	}
	return f.postgresDBSingleton
}

func newFactory() *factory {
	return &factory{}
}
