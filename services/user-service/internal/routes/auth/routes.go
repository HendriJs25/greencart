package auth

import (
	"user-service/internal/handler"
	"user-service/internal/middleware"
	userroutes "user-service/internal/routes/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handlers *handler.Registry, authentication *middleware.Authentication) {
	auth := router.Group("/auth")

	auth.POST("/signin", handlers.AuthHandler.SignIn)

	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	userroutes.Register(authenticated, handlers.UserHandler)
}
