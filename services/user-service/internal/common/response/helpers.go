package response

import (
	"net/http"
	errWrap "user-service/internal/common/error"
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

func BadRequest(c *gin.Context) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusBadRequest,
		Err:  errConstant.ErrBadRequest,
		Gin:  c,
	})
}

func ValidationError(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code:    http.StatusUnprocessableEntity,
		Message: errWrap.ErrValidationResponse(err),
		Err:     err,
		Gin:     c,
	})
}

func Unauthorized(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusUnauthorized,
		Err:  err,
		Gin:  c,
	})
}

func Forbidden(c *gin.Context, err error) {
	HTTPResponse(ParamHTTPResponse{
		Code: http.StatusForbidden,
		Err:  err,
		Gin:  c,
	})
}
