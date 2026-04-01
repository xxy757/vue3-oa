package middleware

import (
	"net/http"
	"strings"

	"oa-saas/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Tenant(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantSlug := c.GetHeader("X-Tenant-Slug")

		if tenantSlug == "" {
			host := c.Request.Host
			host = strings.Split(host, ":")[0]
			parts := strings.Split(host, ".")
			if len(parts) >= 1 && parts[0] != "www" && parts[0] != "admin" && parts[0] != "localhost" && parts[0] != "127" {
				tenantSlug = parts[0]
			}
		}

		if tenantSlug == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "缺少租户标识"})
			c.Abort()
			return
		}

		var tenant model.Tenant
		if err := db.Where("slug = ? AND deleted_at IS NULL", tenantSlug).First(&tenant).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "租户不存在"})
			c.Abort()
			return
		}

		if tenant.Status == "suspended" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "租户已暂停服务"})
			c.Abort()
			return
		}

		if tenant.Status == "cancelled" {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "租户已注销"})
			c.Abort()
			return
		}

		c.Set("tenant_id", tenant.ID)
		c.Set("tenant", tenant)
		c.Next()
	}
}
