package bootstrap

import (
	"os"
	"user-service/internal/config"
	"user-service/internal/logger"

	"github.com/sirupsen/logrus"
)

func SetupLogger(cfg config.App) {
	logger.Log.SetOutput(os.Stdout)

	if cfg.IsDebugMode {
		logger.Log.SetLevel(logrus.DebugLevel)
		logger.Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
		return
	}

	logger.Log.SetLevel(logrus.InfoLevel)
	logger.Log.SetFormatter(&logrus.JSONFormatter{})
}
