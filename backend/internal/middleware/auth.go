package middleware

import (
	"net/http"
	"strings"

	"oa-saas/internal/pkg/jwt"

	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未登录或登录已过期"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "登录已过期"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
