package postgres_test

import (
	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

func connectionDetails() postgres.ConnectionDetails {
	return postgres.ConnectionDetails{
		User:     "postgres",
		Pass:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "postgres",
		SSLMode:  "disable",
	}
}
