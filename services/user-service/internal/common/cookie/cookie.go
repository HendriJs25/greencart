package cookie

import (
	"net/http"
	"time"
	"user-service/internal/constants"

	"github.com/gin-gonic/gin"
)

func SetAccessToken(c *gin.Context, token string, ttl time.Duration) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenName,
		token,
		int(ttl.Seconds()),
		"/",
		"",
		false, // local development
		true,
	)
}

func ClearAccessToken(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		constants.AccessTokenName,
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}
