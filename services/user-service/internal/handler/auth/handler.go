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
		case errors.Is(err, errConstant.ErrUserNotFound), errors.Is(err, errConstant.ErrPasswordIncorrect), errors.Is(err, errConstant.ErrEmailNotVerified):
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

func (h *Handler) SignUp(c *gin.Context) {
	var req request.SignUpRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.authService.SignUp(c.Request.Context(), request.ToSignUpInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrEmailAlreadyExists):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to sign up")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusCreated,
		Message: "アカウントの登録が完了しました",
		Gin:     c,
	})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req request.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.authService.ForgotPassword(c.Request.Context(), req.Email)
	if err != nil {
		if !errors.Is(err, errConstant.ErrUserNotFound) {
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to forgot password")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "ご入力いただいたメールアドレスが登録されている場合、パスワード再設定のご案内をお送りします。メールをご確認ください",
		Gin:     c,
	})
}

func (h *Handler) VerifyAccount(c *gin.Context) {
	var query request.VerifyAccountQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c)
		return
	}

	err := h.authService.VerifyEmail(c.Request.Context(), query.Token)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrUserNotFound), errors.Is(err, errConstant.ErrInvalidToken):
			response.Unauthorized(c, errConstant.ErrInvalidToken)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to verify account")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "アカウントの認証が完了しました",
		Gin:     c,
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var query request.ResetPasswordQuery
	var req request.ResetPasswordRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c)
		return
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.authService.ResetPassword(c.Request.Context(), authservice.ResetPasswordInput{
		Token:       query.Token,
		NewPassword: req.PasswordNew,
	})
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrInvalidToken), errors.Is(err, errConstant.ErrUserNotFound):
			response.Unauthorized(c, errConstant.ErrInvalidToken)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to reset password")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "パスワードの再設定が完了しました",
		Gin:     c,
	})
}
