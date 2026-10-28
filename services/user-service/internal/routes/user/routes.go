package user

import (
	userhandler "user-service/internal/handler/user"
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, userHandler *userhandler.Handler) {
	router.GET("/testauth", userHandler.Auth)
	router.POST("/testcsrf", middleware.CSRFProtection(), userHandler.CSRF)
}
