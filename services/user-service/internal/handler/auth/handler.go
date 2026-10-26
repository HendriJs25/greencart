package auth

import (
	"errors"
	"net/http"
	"time"
	"user-service/internal/common/cookie"
	"user-service/internal/common/response"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/dto/request"
	responsedto "user-service/internal/domain/dto/response"
	"user-service/internal/logger"
	authservice "user-service/internal/services/auth"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	authService authservice.Service
	validate    *customValidator.Validate
}

func NewHandler(authService authservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		authService: authService,
		validate:    validate,
	}
}

func (h *Handler) SignIn(c *gin.Context) {
	var req request.SignInRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	result, err := h.authService.SignIn(c.Request.Context(), request.ToSignInInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrUserNotFound), errors.Is(err, errConstant.ErrPasswordIncorrect):
			response.Unauthorized(c, errConstant.ErrInvalidEmailOrPassword)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to sign in")
			response.InternalServerError(c)
			return
		}
	}

	cookieCreatedAt := time.Now().UTC()
	cookieTTL := result.AccessTokenExpiresAt.Sub(cookieCreatedAt)
	cookie.SetAccessToken(c, result.AccessToken, cookieTTL)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToSignInResponse(result),
		Gin:  c,
	})
}
