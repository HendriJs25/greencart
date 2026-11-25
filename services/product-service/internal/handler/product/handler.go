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
	"strings"

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
	var query request.GetAllProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c)
		return
	}

	if strings.TrimSpace(c.Query("min_price")) == "" {
		query.MinPrice = nil
	}

	if strings.TrimSpace(c.Query("max_price")) == "" {
		query.MaxPrice = nil
	}

	products, metadata, err := h.productService.GetAll(c.Request.Context(), request.ToGetAllProductsQueryInput(query))
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get all products")
		response.InternalServerError(c)
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetAllProductsResponse(products),
		Meta: metadata,
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

func (h *Handler) GetProductHome(c *gin.Context) {
	products, err := h.productService.GetProductHome(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get products for home")
		response.InternalServerError(c)
		return
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: responsedto.ToGetProductHomeResponse(products),
		Gin:  c,
	})
}

func (h *Handler) UploadImage(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	const maxRequestSize = 3 << 20

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxRequestSize)

	file, err := c.FormFile("image")
	if err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			response.UnprocessableEntity(c, errConstant.ErrFileTooBig)
			return
		}
		response.BadRequest(c)
		return
	}

	key, err := h.productService.UploadImage(c.Request.Context(), id, file)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrFileTooBig), errors.Is(err, errConstant.ErrUnsupportedFileExtension):
			response.UnprocessableEntity(c, err)
			return
		case errors.Is(err, errConstant.ErrProductNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).WithField("product_id", id).Error("failed to upload product image")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: key,
		Gin:  c,
	})
}

func (h *Handler) Update(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	err = h.productService.UpdateProduct(c.Request.Context(), id, request.ToUpdateProductInput(req))
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrProductNotFound):
			response.NotFound(c, err)
			return
		case errors.Is(err, errConstant.ErrCategoryNotFound), errors.Is(err, errConstant.ErrAtLeastOneVariantRequired):
			response.UnprocessableEntity(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("product_id", id).WithError(err).Error("failed to update product")
			response.InternalServerError(c)
			return
		}
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	strID := c.Param("id")
	id, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		response.BadRequest(c)
		return
	}

	err = h.productService.DeleteProduct(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrProductNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithField("product_id", id).WithError(err).Error("failed to delete product")
			response.InternalServerError(c)
			return
		}
	}
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}
