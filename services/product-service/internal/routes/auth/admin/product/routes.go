package product

import (
	producthandler "product-service/internal/handler/product"
	"product-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, productHandler *producthandler.Handler) {
	router.POST("/products", middleware.CSRFProtection(), productHandler.Create)
}
