package migration

import (
	"fmt"
	"product-service/internal/config"
	"product-service/internal/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New(cfg *config.Database) (*migrate.Migrate, error) {
	m, err := migrate.New("file://migrations", database.BuildPostgresURL(cfg))
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}
	return m, nil
}
