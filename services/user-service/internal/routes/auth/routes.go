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
	auth.POST("/signup", handlers.AuthHandler.SignUp)
	auth.POST("/forgot-password", handlers.AuthHandler.ForgotPassword)
	auth.POST("/verify-account", handlers.AuthHandler.VerifyAccount)
	auth.PATCH("/reset-password", handlers.AuthHandler.ResetPassword)

	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	userroutes.Register(authenticated, handlers.UserHandler)
}
