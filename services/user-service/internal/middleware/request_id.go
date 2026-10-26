package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"user-service/internal/common/helper"
	"user-service/internal/logger"
	"uuid"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")

		if !validRequestID(requestID) {
			var err error
			requestID, err = generateRequestID()
			if err != nil {
				logger.Log.WithError(err).Warn("failed to generate request id")
			}
		}

		ctx := helper.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

func validRequestID(id string) bool {
	// format Nginx: 32 char hex
	if len(id) == 32 {
		_, err := hex.DecodeString(id)
		return err == nil
	}

	// format UUID standard: 36 char
	if len(id) == 36 {
		_, err := uuid.Parse(id)
		return err == nil
	}
	return false
}

func generateRequestID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return hex.EncodeToString(b[:]), nil
}
