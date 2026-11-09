package user

import (
	"errors"
	"net/http"
	"user-service/internal/common/response"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/dto/request"
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

func (h *Handler) GetMe(c *gin.Context) {
	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	user, err := h.userService.GetMe(c.Request.Context(), identity.UUID)
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).WithField("user_uuid", identity.UUID).Error("get me failed")
		response.InternalServerError(c)
		return
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetMeResponse(user),
		Gin:  c,
	})
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

func (h *Handler) UploadProfileImage(c *gin.Context) {
	const maxRequestSize = 3 << 20

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSize)

	file, err := c.FormFile("photo")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			response.UnprocessableEntity(c, errConstant.ErrFileTooBig)
			return
		}
		response.BadRequest(c)
		return
	}

	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	key, err := h.userService.UploadProfileImage(c.Request.Context(), identity.UUID, file)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrFileTooBig), errors.Is(err, errConstant.ErrUnsupportedFileExtension):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).WithField("user_uuid", identity.UUID).Error("upload profile image failed")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.UploadImageProfileResponse{
			ImageURL: key,
		},
		Gin: c,
	})
}

func (h *Handler) UpdateUserProfileData(c *gin.Context) {
	var req request.UpdateUserDataRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	identity, ok := middleware.RequireIdentity(c)
	if !ok {
		return
	}

	input := request.ToUpdateUserDataInput(req)

	err := h.userService.UpdateUserProfile(c.Request.Context(), identity.UUID, input)
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).WithField("user_uuid", identity.UUID).Error("update user profile failed")
		response.InternalServerError(c)
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}
