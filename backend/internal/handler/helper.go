package handler

import (
	"errors"
	"net/http"

	"oa-saas/internal/pkg/utils"

	"github.com/gin-gonic/gin"
)

func getTenantID(c *gin.Context) uint {
	v, _ := c.Get("tenant_id")
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

func userIDToUint(id interface{}) uint {
	if v, ok := id.(uint); ok {
		return v
	}
	return 0
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *utils.ServiceError
	if errors.As(err, &svcErr) {
		c.JSON(svcErr.Status, gin.H{"code": svcErr.Status, "message": svcErr.Message})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "服务器内部错误"})
}
