package routes

import (
	"product-service/internal/handler"
	"product-service/internal/middleware"
	authroutes "product-service/internal/routes/auth"
	healthroutes "product-service/internal/routes/health"
	productroutes "product-service/internal/routes/product"

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
	authroutes.Register(r.router, r.handlers, r.authentication)
	productroutes.Register(r.router, r.handlers.ProductHandler)
}
