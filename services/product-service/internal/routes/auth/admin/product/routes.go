package product

import (
	producthandler "product-service/internal/handler/product"
	"product-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, productHandler *producthandler.Handler) {
	router.GET("/products", productHandler.GetAll)
	router.GET("/products/:id", productHandler.GetParent)
	router.POST("/products", middleware.CSRFProtection(), productHandler.Create)
	router.POST("/products/:id/image", middleware.CSRFProtection(), productHandler.UploadImage)
	router.PATCH("/products/:id", middleware.CSRFProtection(), productHandler.Update)
	router.DELETE("/products/:id", middleware.CSRFProtection(), productHandler.Delete)
}
