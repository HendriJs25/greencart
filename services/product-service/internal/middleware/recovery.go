package middleware

import (
	"io"
	"product-service/internal/common/response"
	"product-service/internal/logger"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {
		logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"error":  recovered,
			"stack":  string(debug.Stack()),
		})

		response.InternalServerError(c)
		c.Abort()
	})
}
