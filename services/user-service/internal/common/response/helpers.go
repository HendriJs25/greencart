package response

import (
	"net/http"
	errConstant "user-service/internal/constants/error"

	"github.com/gin-gonic/gin"
)

func InternalServerError(c *gin.Context) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusInternalServerError,
		Err:  errConstant.ErrInternalServerError,
		Gin:  c,
	})
}
