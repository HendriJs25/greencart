package bootstrap

import (
	"product-service/internal/config"

	"github.com/gin-gonic/gin"
)

func SetupAppMode(cfg config.App) {
	if cfg.IsDebugMode {
		gin.SetMode(gin.DebugMode)
		return
	}

	gin.SetMode(gin.ReleaseMode)
}
