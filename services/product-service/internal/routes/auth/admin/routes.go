package admin

import (
	"product-service/internal/constants"
	"product-service/internal/handler"
	"product-service/internal/middleware"
	"product-service/internal/routes/auth/admin/category"
	"product-service/internal/routes/auth/admin/product"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, handlers *handler.Registry) {
	admin := router.Group("/admin")
	admin.Use(middleware.RequireRole(constants.Admin))

	category.Register(admin, handlers.CategoryHandler)
	product.Register(admin, handlers.ProductHandler)
}
