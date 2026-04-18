// Package handler 包含所有 HTTP 请求的处理器（Controller 层）。
//
// 每个 handler 对应一个业务模块，负责：
//   - 解析和验证请求参数
//   - 调用数据模型执行业务逻辑
//   - 构造并返回 JSON 响应
//
// 本文件为认证模块处理器，处理用户登录、密码修改和用户信息获取。
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

// ChangePasswordRequest 修改密码请求的参数结构。
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"` // 用户当前的登录密码
	NewPassword string `json:"newPassword" binding:"required,min=6"` // 新密码，最少6位
}

// AuthHandler 认证处理器，持有数据库连接和 JWT 配置。
type AuthHandler struct {
	db        *gorm.DB // 数据库实例
	jwtSecret string   // JWT 签名密钥
	jwtExpire int      // JWT 令牌过期时间（小时）
}

// NewAuthHandler 创建认证处理器实例。
//
// 参数：
//   - db：GORM 数据库实例
//   - cfg：应用配置（从中提取 JWT 密钥和过期时间）
//
// 返回值：
//   - *AuthHandler：初始化后的认证处理器
func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		db:        db,
		jwtSecret: cfg.JWT.Secret,
		jwtExpire: cfg.JWT.ExpireHours,
	}
}

// LoginRequest 用户登录请求的参数结构。
type LoginRequest struct {
	Username string `json:"username" binding:"required"` // 用户名
	Password string `json:"password" binding:"required"` // 登录密码
}

// Login 处理用户登录请求。
// 路由：POST /api/v1/auth/login
//
// 功能：验证用户凭据，生成 JWT 令牌，返回用户信息和企业信息。
//
// 参数：
//   - c：Gin 上下文（从中获取 tenant_id，以及请求体中的用户名和密码）
//
// 响应数据包含：
//   - token：JWT 访问令牌
//   - user：用户基本信息（ID、用户名、昵称、邮箱、手机、头像、部门ID、角色ID）
//   - tenant：所属企业信息（ID、名称、标识、状态）
//   - tenant.plan：企业当前套餐信息
//
// 执行步骤：
//  1. 绑定并验证请求参数（用户名、密码）
//  2. 从上下文获取租户 ID，在租户范围内查找用户
//  3. 检查用户状态（是否被禁用）
//  4. 使用 bcrypt 验证密码
//  5. 生成 JWT 令牌
//  6. 查询企业信息和套餐信息
//  7. 组装并返回登录响应
func (h *AuthHandler) Login(c *gin.Context) {
	// 步骤1：绑定并验证请求参数
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "用户名和密码不能为空"})
		return
	}

	// 步骤2：在租户范围内查找用户
	tenantID, _ := c.Get("tenant_id")
	var user model.User
	if err := h.db.Where("username = ? AND tenant_id = ?", req.Username, tenantID.(uint)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 步骤3：检查用户状态
	if user.Status != 1 {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "用户已被禁用"})
		return
	}

	// 步骤4：验证密码（bcrypt 哈希比对）
	if !utils.CheckPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "用户名或密码错误"})
		return
	}

	// 步骤5：生成 JWT 令牌
	token, err := jwt.GenerateToken(user.ID, user.TenantID, h.jwtSecret, h.jwtExpire)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成token失败"})
		return
	}

	// 步骤6：查询企业和套餐信息
	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	var plan model.Plan
	h.db.First(&plan, tenant.PlanID)

	// 步骤7：返回登录响应
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

// ChangePassword 处理用户修改密码请求。
// 路由：PUT /api/v1/auth/password（需认证）
//
// 功能：验证旧密码正确后，将密码更新为新密码的 bcrypt 哈希值。
//
// 参数：
//   - c：Gin 上下文（包含 user_id 和 tenant_id）
//
// 执行步骤：
//  1. 绑定请求参数（旧密码、新密码）
//  2. 从数据库查询当前用户
//  3. 验证旧密码是否正确
//  4. 对新密码进行 bcrypt 哈希
//  5. 更新数据库中的密码字段
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")

	// 步骤1：绑定请求参数
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	// 步骤2：查询当前用户
	var user model.User
	if err := h.db.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	// 步骤3：验证旧密码
	if !utils.CheckPassword(req.OldPassword, user.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "原密码错误"})
		return
	}

	// 步骤4：对新密码进行哈希
	hashedPwd, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	// 步骤5：更新密码
	h.db.Model(&user).Update("password", hashedPwd)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "密码修改成功"})
}

// GetInfo 获取当前登录用户的详细信息。
// 路由：GET /api/v1/auth/info（需认证）
//
// 功能：返回当前用户的完整信息，包括部门名称、角色名称、权限列表和企业信息。
//
// 参数：
//   - c：Gin 上下文（包含 user_id 和 tenant_id）
//
// 响应数据包含：
//   - 用户基本信息（ID、用户名、昵称、邮箱、手机、头像等）
//   - deptName：所属部门名称
//   - roleName：角色名称
//   - permissions：权限列表（来自角色配置）
//   - tenant：企业信息及当前套餐
//
// 执行步骤：
//  1. 从上下文获取 user_id 和 tenant_id，查询用户记录
//  2. 若用户有所属部门，查询部门名称
//  3. 若用户有角色，查询角色名称和权限列表
//  4. 从上下文获取租户信息，查询套餐详情
//  5. 组装并返回完整的用户信息
func (h *AuthHandler) GetInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tenantID, _ := c.Get("tenant_id")

	// 步骤1：查询用户基本信息
	var user model.User
	if err := h.db.Where("id = ? AND tenant_id = ?", userID, tenantID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	// 步骤2：查询部门名称
	var deptName string
	if user.DeptID != nil {
		var dept model.Department
		if err := h.db.Where("id = ? AND tenant_id = ?", *user.DeptID, tenantID).First(&dept).Error; err == nil {
			deptName = dept.Name
		}
	}

	// 步骤3：查询角色名称和权限列表
	var roleName string
	var permissions model.StringArray
	if user.RoleID != nil {
		var role model.Role
		if err := h.db.Where("id = ? AND tenant_id = ?", *user.RoleID, tenantID).First(&role).Error; err == nil {
			roleName = role.Name
			permissions = role.Permissions
		}
	}

	// 步骤4：获取租户和套餐信息
	tenantObj, _ := c.Get("tenant")
	tenant := tenantObj.(model.Tenant)
	var plan model.Plan
	h.db.First(&plan, tenant.PlanID)

	// 步骤5：返回完整用户信息
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
