package cmd

import (
	"fmt"
	"user-service/internal/database"
	"user-service/internal/logger"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start user service",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func runServer() error {
	postgresDB, err := database.NewPostgres(cfg)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}

	defer func() {
		if err := postgresDB.Close(); err != nil {
			logger.Log.WithError(err).Warn("failed to close postgres connection")
		}
	}()

	logger.Log.WithFields(logrus.Fields{
		"host": cfg.Database.Host,
		"port": cfg.Database.Port,
		"name": cfg.Database.Name,
	}).Info("postgres connection established")

	logger.Log.WithFields(logrus.Fields{
		"environment": cfg.App.Env,
		"address":     cfg.ServerAddress(),
	}).Info("starting user service")

	return nil
}
