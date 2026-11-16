package cmd

import (
	"fmt"
	"product-service/internal/bootstrap"
	"product-service/internal/config"

	"github.com/spf13/cobra"
)

var cfg *config.Config

var rootCmd = &cobra.Command{
	Use:   "product-service",
	Short: "product-service",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		cfg, err = config.Load()
		if err != nil {
			return fmt.Errorf("load config failed: %w", err)
		}

		bootstrap.SetupLogger(cfg.App)
		return nil
	},
}

func Execute() error {
	registerCommands()
	return rootCmd.Execute()
}

func registerCommands() {
	rootCmd.AddCommand(
		serveCmd,
		migrateCmd,
	)

	migrateCmd.AddCommand(
		migrateUpCmd,
		migrateDownCmd,
		migrateForceCmd,
		migrateVersionCmd,
	)
}
