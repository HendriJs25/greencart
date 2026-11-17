package category

import (
	categoryhandler "product-service/internal/handler/category"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, categoryHandler *categoryhandler.Handler) {
	router.GET("/categories", categoryHandler.GetAll)
}
