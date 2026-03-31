package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{Code: 200, Data: data})
}

func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, APIResponse{Code: httpStatus, Message: message})
}
