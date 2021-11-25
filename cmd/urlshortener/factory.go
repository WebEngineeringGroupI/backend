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
	shortURLRepositorySingleton url.ShortURLRepository
	urlValidatorSingleton       url.Validator
}

func (f *factory) NewHTTPRouter() gohttp.Handler {
	return http.NewRouter(f.httpConfig())
}

func (f *factory) httpConfig() http.Config {
	return http.Config{
		BaseDomain:         f.baseDomain(),
		ShortURLRepository: f.shortURLRepository(),
		URLValidator:       f.urlValidator(),
		CustomMetrics:      f.customMetrics(),
	}
}

func (f *factory) customMetrics() url.Metrics {
	return metrics.NewPrometheusMetrics()
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
		db, err := postgres.NewDB(f.postgresConnectionDetails())
		if err != nil {
			log.Fatalf("unable to create the database connection: %s", err)
		}
		f.shortURLRepositorySingleton = db
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

func (f *factory) NewGRPCServer() *gogrpc.Server {
	return grpc.NewServer(grpc.Config{
		BaseDomain:         f.baseDomain(),
		ShortURLRepository: f.shortURLRepository(),
		URLValidator:       f.urlValidator(),
	})
}

func newFactory() *factory {
	return &factory{}
}
