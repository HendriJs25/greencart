package middleware

import (
	"io"

	"github.com/gin-gonic/gin"
)

func HandlePanic() gin.HandlerFunc {
	return gin.CustomRecoveryWithWriter(io.Discard, func(c *gin.Context, recovered any) {

	})
}
