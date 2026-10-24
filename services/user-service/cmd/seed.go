package cmd

import (
	"context"
	"fmt"
	"time"
	"user-service/internal/database"
	"user-service/internal/database/seeders"
	"user-service/internal/logger"

	"github.com/spf13/cobra"
)

const seedTimeOut = 30 * time.Second

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "run seeders",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSeeder()
	},
}

func runSeeder() error {
	postgresDB, err := database.NewPostgres(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer func() {
		if err := postgresDB.Close(); err != nil {
			logger.Log.WithError(err).Warn("failed to close postgres connection")
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), seedTimeOut)
	defer cancel()

	err = seeders.Run(ctx, postgresDB.DB, cfg.Seed)

	logger.Log.Info("starting seeders")

	if err != nil {
		return err
	}

	logger.Log.Info("seeders finished")
	return nil
}
