package customer

import (
	"errors"
	"net/http"
	"strconv"
	"user-service/internal/common/response"
	errConstant "user-service/internal/constants/error"
	"user-service/internal/domain/dto/request"
	responsedto "user-service/internal/domain/dto/response"
	"user-service/internal/logger"
	customerservice "user-service/internal/services/customer"

	"github.com/gin-gonic/gin"
	customValidator "github.com/go-playground/validator/v10"
)

type Handler struct {
	customerService customerservice.Service
	validate        *customValidator.Validate
}

func NewHandler(customerService customerservice.Service, validate *customValidator.Validate) *Handler {
	return &Handler{
		customerService: customerService,
		validate:        validate,
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req request.CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.ValidationError(c, err)
		return
	}

	input := request.ToCreateCustomerInput(req)

	err := h.customerService.Create(c.Request.Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrEmailAlreadyExists):
			response.Conflict(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).Error("failed to create customer")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code:    http.StatusOK,
		Message: "顧客を登録しました",
		Gin:     c,
	})
}

func (h *Handler) GetAll(c *gin.Context) {
	var query request.GetAllCustomersQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c)
		return
	}

	customers, metadata, err := h.customerService.GetAll(c.Request.Context(), request.ToGetAllCustomerQueryInput(query))
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get all customer")
		response.InternalServerError(c)
		return
	}

	result := responsedto.ToGetAllCustomerResponse(customers)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Meta: metadata,
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

	customer, err := h.customerService.GetDetail(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrUserNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).WithField("customer_id", id).Error("failed to get customer detail")
			response.InternalServerError(c)
			return
		}
	}

	result := responsedto.ToGetDetailCustomerResponse(customer)
	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
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

	err = h.customerService.Delete(c.Request.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, errConstant.ErrUserNotFound):
			response.NotFound(c, err)
			return
		default:
			logger.WithContext(c.Request.Context()).WithError(err).WithField("customer_id", strID).Error("failed to delete customer")
			response.InternalServerError(c)
			return
		}
	}

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Gin:  c,
	})
}
