package product

import (
	"errors"
	"net/http"
	"product-service/internal/common/response"
	errConstant "product-service/internal/constants/error"
	"product-service/internal/domain/dto/request"
	responsedto "product-service/internal/domain/dto/response"
	"product-service/internal/logger"
	productservice "product-service/internal/services/product"
	"strconv"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	productService productservice.Service
	validate       *customValidator.Validate
}

func NewHandler(productService productservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		productService: productService,
		validate:       validate,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err := h.productService.Create(c.Request.Context(), request.ToCreateProductInput(req))
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

func (h *Handler) GetAll(c *gin.Context) {
	products, err := h.productService.GetAll(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get all products")
		response.InternalServerError(c)
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetAllProductsResponse(products),
		Gin:  c,
	})
}

func (h *Handler) GetParent(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	product, err := h.productService.GetParent(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrProductNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("product_id", id).WithError(err).Error("failed to get product")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetProductDetailResponse(product),
		Gin:  c,
	})
}
