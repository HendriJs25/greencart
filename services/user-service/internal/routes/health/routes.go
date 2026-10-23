package health

import (
	healthhandler "user-service/internal/handler/health"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, healthHandler *healthhandler.Handler) {
	router.GET("/healthcheck", healthHandler.HealthCheck)
}
