package middleware

import (
	"user-service/internal/common/response"
	errConstant "user-service/internal/constants/error"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	allowedRoleSet := make(map[string]struct{}, len(allowedRoles))

	for _, role := range allowedRoles {
		allowedRoleSet[role] = struct{}{}
	}

	return func(c *gin.Context) {
		identity, ok := RequireIdentity(c)
		if !ok {
			return
		}

		allowed := 0

		for _, roleName := range identity.RoleName {
			if _, ok := allowedRoleSet[roleName]; ok {
				allowed++
			}
		}

		if allowed == 0 {
			response.Forbidden(c, errConstant.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
