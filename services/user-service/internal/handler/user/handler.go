package user

import "github.com/gin-gonic/gin"

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Auth(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}

func (h *Handler) CSRF(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
