package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/service"

	"github.com/gin-gonic/gin"
)

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名和密码不能为空"})
		return
	}
	tenantID, _ := c.Get("tenant_id")
	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	result, err := h.authService.Login(req.Username, tenantID.(uint), req.Password, tenant)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.authService.ChangePassword(userID.(uint), tenantID.(uint), req.OldPassword, req.NewPassword); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码修改成功"})
}

func (h *AuthHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	result, err := h.authService.GetInfo(userID.(uint), tenantID.(uint), tenant)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}
