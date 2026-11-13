package customer

import (
	customerhandler "user-service/internal/handler/customer"
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, customerHandler *customerhandler.Handler) {
	router.POST("/customers", middleware.CSRFProtection(), customerHandler.Create)
	router.GET("/customers", customerHandler.GetAll)
	router.GET("/customers/:id", customerHandler.GetDetail)
	router.PATCH("/customers/:id", middleware.CSRFProtection(), customerHandler.Update)
	router.DELETE("/customers/:id", middleware.CSRFProtection(), customerHandler.Delete)
}
