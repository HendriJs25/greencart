package middleware

import (
	"time"
	"user-service/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		path := c.Request.URL.Path
		query := c.Request.URL.Query()
		query.Del("token")

		if encoded := query.Encode(); encoded != "" {
			path += "?" + encoded
		}

		entry := logger.WithContext(c.Request.Context()).WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       path,
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  c.ClientIP(),
		})

		entry.Info("request completed")
	}
}
