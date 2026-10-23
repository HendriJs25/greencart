package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"user-service/internal/common/requestcontext"
	"uuid"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-Id")

		if !validRequestID(requestID) {
			requestID = generateRequestID()
		}

		ctx := requestcontext.WithRequestID(c.Request.Context(), requestID)
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

func generateRequestID() string {
	var b [16]byte
	rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
