package postgres_test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/WebEngineeringGroupI/backend/pkg/infrastructure/database/postgres"
)

func connectionDetails() *postgres.ConnectionDetails {
	return &postgres.ConnectionDetails{
		User:     "postgres",
		Pass:     "postgres",
		Host:     "localhost",
		Port:     5432,
		Database: "postgres",
		SSLMode:  "disable",
	}
}

func randomHash() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Int())[0:7]
}
