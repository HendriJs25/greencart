package user

import (
	userhandler "user-service/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, userHandler *userhandler.Handler) {
	router.GET("/profile", userHandler.GetProfile)
}
