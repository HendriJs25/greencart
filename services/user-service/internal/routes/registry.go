package routes

import (
	"user-service/internal/handler"
	healthroutes "user-service/internal/routes/health"

	"github.com/gin-gonic/gin"
)

type Registry struct {
	router   *gin.RouterGroup
	handlers *handler.Registry
}

func NewRegistry(router *gin.RouterGroup, handlers *handler.Registry) *Registry {
	return &Registry{
		router:   router,
		handlers: handlers,
	}
}

func (r *Registry) RegisterRoutes() {
	healthroutes.Register(r.router, r.handlers.HealthHandler)
}
