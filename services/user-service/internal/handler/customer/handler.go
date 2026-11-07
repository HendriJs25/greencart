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
)

type Handler struct {
	customerService customerservice.Service
}

func NewHandler(customerService customerservice.Service) *Handler {
	return &Handler{
		customerService: customerService,
	}
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
			logger.WithContext(c.Request.Context()).WithError(err).WithField("customer_id", customer.ID).Error("failed to get customer detail")
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
