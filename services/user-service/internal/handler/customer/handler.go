package customer

import (
	"net/http"
	"user-service/internal/common/response"
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
	customers, err := h.customerService.GetAll(c.Request.Context())
	if err != nil {
		logger.WithContext(c.Request.Context()).WithError(err).Error("failed to get all customer")
		response.InternalServerError(c)
		return
	}

	result := responsedto.ToGetAllCustomerResponse(customers)

	response.HTTPResponse(response.ParamHTTPResponse{
		Code: http.StatusOK,
		Data: result,
		Gin:  c,
	})
}
