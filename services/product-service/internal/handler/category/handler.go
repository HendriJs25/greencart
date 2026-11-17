package category

import (
	"net/http"
	"product-service/internal/common/response"
	responsedto "product-service/internal/domain/dto/response"
	"product-service/internal/logger"
	categoryservice "product-service/internal/services/category"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	categoryService categoryservice.Service
}

func NewHandler(categoryService categoryservice.Service) *Handler {
	return &Handler{
		categoryService: categoryService,
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
