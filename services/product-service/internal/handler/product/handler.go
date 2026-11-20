package product

import (
	"errors"
	"net/http"
	"product-service/internal/common/response"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/dto/request"
	"product-service/internal/logger"
	productservice "product-service/internal/services/product"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	ProductService productservice.Service
	Validate       *customValidator.Validate
}

func NewHandler(productService productservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		ProductService: productService,
		Validate:       validate,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.ProductService.Create(c.Request.Context(), request.ToCreateProductInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrCategoryNotFound):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to create product")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusCreated,
		Gin:  c,
	})
}
