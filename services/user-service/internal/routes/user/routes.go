package user

import (
	userhandler "user-service/internal/handler/user"
	"user-service/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.RouterGroup, userHandler *userhandler.Handler) {
	router.GET("/profile", userHandler.GetProfile)
	router.PATCH("/profile/image-upload", middleware.CSRFProtection(), userHandler.UploadProfileImage)
}
