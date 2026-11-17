package routes

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/routes/auth"
	healthroutes "product-service/internal/routes/health"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	router         *gin.RouterGroup
	handlers       *handler.Registry
	authentication *middleware.Authentication
}

func NewRegistry(router *gin.RouterGroup, handlers *handler.Registry, authentication *middleware.Authentication) *Registry {
	return &Registry{
		router:         router,
		handlers:       handlers,
		authentication: authentication,
	}
}

func (r *Registry) RegisterRoutes() {
	healthroutes.Register(r.router, r.handlers.HealthHandler)
	auth.Register(r.router, r.handlers, r.authentication)
}
