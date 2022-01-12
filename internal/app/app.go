package app

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

func PostgresConnectionDetails() *postgres.ConnectionDetails {
	dbPort, err := strconv.Atoi(mandatoryEnvVarValue("DB_PORT"))
	if err != nil {
		log.Fatalf("unable to parse DB_PORT, make sure it is defined and is a valid integer")
	}

	return &postgres.ConnectionDetails{
		User:     mandatoryEnvVarValue("DB_USER"),
		Pass:     mandatoryEnvVarValue("DB_PASS"),
		Host:     mandatoryEnvVarValue("DB_HOST"),
		Port:     dbPort,
		Database: mandatoryEnvVarValue("DB_NAME"),
		SSLMode:  mandatoryEnvVarValue("DB_SSL_MODE"),
	}
}

func RabbitMQConnectionString() string {
	rabbitMQPort, err := strconv.Atoi(optionalEnvVarValue("RABBITMQ_PORT", "5672"))
	if err != nil {
		log.Fatalf("unable to parse RABBITMQ_PORT as int, make sure it has a valid value")
	}

	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		mandatoryEnvVarValue("RABBITMQ_USER"),
		mandatoryEnvVarValue("RABBITMQ_PASS"),
		mandatoryEnvVarValue("RABBITMQ_HOST"),
		rabbitMQPort)
}

func SafeBrowsingAPIKey() string {
	return mandatoryEnvVarValue("SAFE_BROWSING_API_KEY")
}

func mandatoryEnvVarValue(variable string) string {
	value, isSet := os.LookupEnv(variable)
	if !isSet {
		log.Fatalf("mandatory %s env var is not set", variable)
	}
	return value
}

func optionalEnvVarValue(variable string, defaultValue string) string {
	value, isSet := os.LookupEnv(variable)
	if !isSet {
		return defaultValue
	}
	return value
}
