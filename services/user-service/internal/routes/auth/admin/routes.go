package admin

import (
	"user-service/internal/constants"
	customerhandler "user-service/internal/handler/customer"
	"user-service/internal/middleware"
	customerroutes "user-service/internal/routes/auth/admin/customer"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, customerHandler *customerhandler.Handler) {
	admin := router.Group("/admin")
	admin.Use(middleware.RequireRole(constants.Admin))

	customerroutes.Register(admin, customerHandler)
}
