package auth

import (
	authhandler "user-service/internal/handler/auth"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, authHandler *authhandler.Handler) {
	auth := router.Group("/auth")

	auth.POST("/signin", authHandler.SignIn)
}
