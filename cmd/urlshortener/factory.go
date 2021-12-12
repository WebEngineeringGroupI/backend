package main

import (
	"log"
	gohttp "net/http"
	"os"
	"strconv"

	gogrpc "google.golang.org/grpc"

	"github.com/WebEngineeringGroupI/backend/pkg/application/grpc"
	"github.com/WebEngineeringGroupI/backend/pkg/application/http"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/event"
	"github.com/WebEngineeringGroupI/backend/pkg/domain/url"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/clock"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/metrics"
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/uuid"
)

type factory struct {
	metricsSingleton    url.Metrics
	postgresDBSingleton *postgres.DBSession
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
		CustomMetrics:              f.customMetrics(),
		ShortURLRepository:         f.newPostgresDB(),
		LoadBalancedURLsRepository: f.newPostgresDB(),
		EventEmitter:               f.eventEmitter(),
	}
}

func (f *factory) grpcConfig() grpc.Config {
	return grpc.Config{
		BaseDomain:                 f.baseDomain(),
		ShortURLRepository:         f.newPostgresDB(),
		CustomMetrics:              f.customMetrics(),
		LoadBalancedURLsRepository: f.newPostgresDB(),
		EventEmitter:               f.eventEmitter(),
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

func (f *factory) newPostgresDB() *postgres.DBSession {
	if f.postgresDBSingleton == nil {
		db, err := postgres.NewDB(f.postgresConnectionDetails())
		if err != nil {
			log.Fatalf("unable to create the database connection: %s", err)
		}
		f.postgresDBSingleton = db.Session()
	}
	return f.postgresDBSingleton
}

func (f *factory) eventEmitter() event.Emitter {
	return event.NewEmitter(f.newPostgresDB(), clock.NewFromSystem(), uuid.NewGenerator())
}

func newFactory() *factory {
	return &factory{}
}
