package user

import (
	"net/http"
	"user-service/internal/common/response"
	responsedto "user-service/internal/domain/dto/response"
	"user-service/internal/logger"
	"user-service/internal/middleware"
	userservice "user-service/internal/services/user"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	userService userservice.Service
	validate    *customValidator.Validate
}

func NewHandler(userService userservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		userService: userService,
		validate:    validate,
	}
}

func (h *Handler) GetProfile(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), identity.UUID)
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).WithField("user_uuid", identity.UUID).Error("get user profile failed")
		response.InternalServerError(c)
		return
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetProfileResponse(user),
		Gin:  c,
	})
}
