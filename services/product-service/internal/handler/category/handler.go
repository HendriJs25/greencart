package category

import (
	"errors"
	"net/http"
	"product-service/internal/common/response"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/dto/request"
	responsedto "product-service/internal/domain/dto/response"
	"product-service/internal/logger"
	categoryservice "product-service/internal/services/category"
	"strconv"

	"github.com/gin-gonic/gin"
	customvalidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	categoryService categoryservice.Service
	validate        *customvalidator.Validate
}

func NewHandler(categoryService categoryservice.Service, validate *customvalidator.Validate) *Handler {
	return &Handler{
		categoryService: categoryService,
		validate:        validate,
	}
}

func (h *Handler) GetAll(c *gin.Context) {
	categories, err := h.categoryService.GetAll(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get categories")
		response.InternalServerError(c)
		return
	}

	result := responsedto.ToGetAllCategoriesResponse(categories)
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (h *Handler) GetDetail(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	category, err := h.categoryService.GetByID(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("category_id", id).WithError(err).Error("failed to get category detail")
			response.InternalServerError(c)
			return
		}
	}

	result := responsedto.ToGetDetailCategoryResponse(category)
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}

func (h *Handler) GetAllSlugForOptions(c *gin.Context) {
	slugs, err := h.categoryService.GetAllSlugForOptions(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get categories slugs")
		response.InternalServerError(c)
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: slugs,
		Gin:  c,
	})
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.categoryService.Create(c.Request.Context(), request.ToCreateCategoryInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrCategoryAlreadyExists):
			response.Conflict(c, err)
			return
		case errors.Is(err, errConstant.ErrNotFound):
			response.UnprocessableEntity(c, errConstant.ErrParentCategoryMissing)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to create category")
			response.InternalServerError(c)
			return
		}
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}

func (h *Handler) UploadIconImage(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	const maxRequestSize = 3 << 20

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSize)

	file, err := c.FormFile("icon")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			response.UnprocessableEntity(c, errConstant.ErrFileTooBig)
			return
		}
		response.BadRequest(c)
		return
	}

	key, err := h.categoryService.UploadIconImage(c.Request.Context(), id, file)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrFileTooBig), errors.Is(err, errConstant.ErrUnsupportedFileExtension):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).WithField("category_id", id).Error("upload icon image failed")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.UploadImageIconResponse{
			IconURL: key,
		},
		Gin: c,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = h.categoryService.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrCategoryHasChildren):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("category_id", id).WithError(err).Error("failed to delete category")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}
