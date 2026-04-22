package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type ServiceError struct {
	Status  int
	Message string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("service error %d: %s", e.Status, e.Message)
}

func ErrBadRequest(msg string) error {
	return &ServiceError{Status: http.StatusBadRequest, Message: msg}
}

func ErrUnauthorized(msg string) error {
	return &ServiceError{Status: http.StatusUnauthorized, Message: msg}
}

func ErrForbidden(msg string) error {
	return &ServiceError{Status: http.StatusForbidden, Message: msg}
}

func ErrNotFound(msg string) error {
	return &ServiceError{Status: http.StatusNotFound, Message: msg}
}

func ErrConflict(msg string) error {
	return &ServiceError{Status: http.StatusConflict, Message: msg}
}

func ErrInternal(msg string) error {
	return &ServiceError{Status: http.StatusInternalServerError, Message: msg}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, APIResponse{
		Code: 200,
		Data: data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Code:    code,
		Message: message,
	})
}
