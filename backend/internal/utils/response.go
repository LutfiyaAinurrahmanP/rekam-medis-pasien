package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data interface{}) {
	ctx.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err interface{}) {
	ctx.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func ValidationErrorResponse(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: "Validation error",
		Error:   err.Error(),
	})
}
