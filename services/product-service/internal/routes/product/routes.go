package product

import (
	producthandler "product-service/internal/handler/product"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, productHandler *producthandler.Handler) {
	router.GET("/home", productHandler.GetProductHome)
	router.GET("/shop", productHandler.GetProductShop)
}
