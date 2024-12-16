package apihelpers

import (
	"auth-rest-api/constants"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIRes struct {
	Status    bool        `json:"status"`
	Message   string      `json:"message"`
	ErrorCode string      `json:"errorcode"`
	Data      interface{} `json:"data"`
}

func ErrorMessageController(c *gin.Context, ErrorCode string) {
	var apiRes APIRes
	apiRes.ErrorCode = ErrorCode
	apiRes.Message = constants.ErrorCodeMap[ErrorCode]
	apiRes.Status = false
	CustomResponse(c, http.StatusBadRequest, apiRes)
}

func SendInternalServerError() (int, APIRes) {
	var apiRes APIRes
	apiRes.Status = false
	apiRes.Message = constants.ErrorCodeMap[constants.InternalServerError]
	apiRes.ErrorCode = constants.InternalServerError
	return http.StatusInternalServerError, apiRes
}

func SendErrorResponse(status bool, errorCode string, httpCode int) (int, APIRes) {
	var apiRes APIRes
	apiRes.Status = status
	apiRes.Message = constants.ErrorCodeMap[errorCode]
	apiRes.ErrorCode = errorCode
	return httpCode, apiRes
}

func CustomResponse(ctx *gin.Context, code int, data interface{}) {

	ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Add("Vary", "Origin")
	ctx.Writer.Header().Add("Vary", "Access-Control-Request-Method")
	ctx.Writer.Header().Add("Vary", "Access-Control-Request-Headers")
	ctx.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Accept, token")
	ctx.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	ctx.JSON(code, data)
}
