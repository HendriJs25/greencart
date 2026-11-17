package auth

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	adminroutes "product-service/internal/routes/auth/admin"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handlers *handler.Registry, authentication *middleware.Authentication) {
	authenticated := router.Group("")
	authenticated.Use(authentication.Handle())

	adminroutes.Register(authenticated, handlers)
}
