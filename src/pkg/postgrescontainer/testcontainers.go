// Package postgrescontainer is required for creating test postgres container with testcontainers package ["github.com/testcontainers/testcontainers-go"].
package postgrescontainer

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

// Config struct is required for container configuration.
type Config struct {
	User     string
	Password string
	DBName   string
}

// New creates a new postgres container and returns an instance of *postgres.PostgresContainer.
func New(config *Config) (*postgres.PostgresContainer, error) {
	postgresContainer, err := postgres.Run(context.Background(),
		"postgres:16.3-alpine3.20",
		postgres.WithDatabase(config.DBName),
		postgres.WithUsername(config.User),
		postgres.WithPassword(config.Password),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("run postgres container: %w", err)
	}

	return postgresContainer, nil
}
