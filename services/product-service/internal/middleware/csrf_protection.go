package middleware

import (
	"product-service/internal/common/helper"
	"product-service/internal/common/response"
	"product-service/internal/constants"
	errConstant "product-service/internal/constants/error"

	"github.com/gin-gonic/gin"
)

func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		identity, ok := RequireIdentity(c)
		if !ok {
			return
		}

		csrfToken := c.GetHeader(constants.CSRFToken)
		if csrfToken == "" {
			response.Unauthorized(c, errConstant.ErrInvalidCSRFToken)
			c.Abort()
			return
		}

		hashedCSRFToken := helper.HashToken(csrfToken)
		if identity.CSRFTokenHash != hashedCSRFToken {
			response.Forbidden(c, errConstant.ErrInvalidCSRFToken)
			c.Abort()
			return
		}

		c.Next()
	}
}
