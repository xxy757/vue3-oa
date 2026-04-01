package handler

import (
	"net/http"

	"oa-saas/internal/config"
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/jwt"
	"oa-saas/internal/pkg/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	jwtSecret string
	jwtExpire int
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: cfg.JWT.Secret,
		jwtExpire: cfg.JWT.ExpireHours,
	}
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
	var user model.User
	if err := h.db.Where("username = ? AND tenant_id = ?", req.Username, tenantID.(uint)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "用户已被禁用"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.TenantID, h.jwtSecret, h.jwtExpire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	var plan model.Plan
	h.db.First(&plan, tenant.PlanID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"nickname": user.Nickname,
				"email":    user.Email,
				"phone":    user.Phone,
				"avatar":   user.Avatar,
				"deptId":   user.DeptID,
				"roleId":   user.RoleID,
			},
			"tenant": gin.H{
				"id":     tenant.ID,
				"name":   tenant.Name,
				"slug":   tenant.Slug,
				"status": tenant.Status,
				"plan": gin.H{
					"id":       plan.ID,
					"name":     plan.Name,
					"code":     plan.Code,
					"features": plan.Features,
					"maxUsers": plan.MaxUsers,
				},
			},
		},
	})
}

func (h *AuthHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")
	var user model.User
	if err := h.db.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	var deptName string
	if user.DeptID != nil {
		var dept model.Department
		if err := h.db.Where("id = ? AND tenant_id = ?", *user.DeptID, tenantID).First(&dept).Error; err == nil {
			deptName = dept.Name
		}
	}

	var roleName string
	var permissions model.StringArray
	if user.RoleID != nil {
		var role model.Role
		if err := h.db.Where("id = ? AND tenant_id = ?", *user.RoleID, tenantID).First(&role).Error; err == nil {
			roleName = role.Name
			permissions = role.Permissions
		}
	}

	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	var plan model.Plan
	h.db.First(&plan, tenant.PlanID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":          user.ID,
			"username":    user.Username,
			"nickname":    user.Nickname,
			"email":       user.Email,
			"phone":       user.Phone,
			"avatar":      user.Avatar,
			"deptId":      user.DeptID,
			"roleId":      user.RoleID,
			"deptName":    deptName,
			"roleName":    roleName,
			"permissions": permissions,
			"createTime":  user.CreatedAt,
			"tenant": gin.H{
				"id":     tenant.ID,
				"name":   tenant.Name,
				"slug":   tenant.Slug,
				"status": tenant.Status,
				"plan": gin.H{
					"id":       plan.ID,
					"name":     plan.Name,
					"code":     plan.Code,
					"features": plan.Features,
					"maxUsers": plan.MaxUsers,
				},
			},
		},
	})
}
