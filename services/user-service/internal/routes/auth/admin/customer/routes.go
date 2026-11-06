package customer

import (
	customerhandler "user-service/internal/handler/customer"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, customerHandler *customerhandler.Handler) {
	router.GET("/checkadmin", customerHandler.CheckAdmin)
}
