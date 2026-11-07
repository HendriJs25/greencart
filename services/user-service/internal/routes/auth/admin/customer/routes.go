package customer

import (
	customerhandler "user-service/internal/handler/customer"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, customerHandler *customerhandler.Handler) {
	router.GET("/customers", customerHandler.GetAll)
	router.GET("/customers/:id", customerHandler.GetDetail)
	router.DELETE("/customers/:id", customerHandler.Delete)
}
