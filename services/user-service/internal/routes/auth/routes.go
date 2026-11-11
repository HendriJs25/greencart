package auth

import (
	"user-service/internal/handler"
	"user-service/internal/middleware"
	adminroutes "user-service/internal/routes/auth/admin"
	userroutes "user-service/internal/routes/auth/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handlers *handler.Registry, authentication *middleware.Authentication) {
	auth := router.Group("/auth")

	auth.POST("/signin", handlers.AuthHandler.SignIn)
	auth.POST("/signup", handlers.AuthHandler.SignUp)
	auth.POST("/signout", authentication.Handle(), middleware.CSRFProtection(), handlers.AuthHandler.SignOut)
	auth.POST("/forgot-password", handlers.AuthHandler.ForgotPassword)
	auth.POST("/verify-account", handlers.AuthHandler.VerifyAccount)
	auth.PATCH("/reset-password", handlers.AuthHandler.ResetPassword)

	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	userroutes.Register(authenticated, handlers.UserHandler)
	adminroutes.Register(authenticated, handlers.CustomerHandler)
}
