package category

import (
	categoryhandler "product-service/internal/handler/category"
	"product-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, categoryHandler *categoryhandler.Handler) {
	router.GET("/categories", categoryHandler.GetAll)
	router.POST("/categories", middleware.CSRFProtection(), categoryHandler.Create)
}
