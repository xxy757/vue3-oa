package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func getTenantID(c *gin.Context) uint {
	v, _ := c.Get("tenant_id")
	if id, ok := v.(uint); ok {
		return id
	}
	return 0
}

func (h *UserHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.DefaultQuery("keyword", "")
	tid := getTenantID(c)

	var users []model.User
	query := h.db.Model(&model.User{}).Where("tenant_id = ?", tid)

	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	query.Count(&total)

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取用户列表失败"})
		return
	}

	type UserItem struct {
		ID         uint   `json:"id"`
		Username   string `json:"username"`
		Nickname   string `json:"nickname"`
		Email      string `json:"email"`
		Phone      string `json:"phone"`
		Avatar     string `json:"avatar"`
		DeptID     *uint  `json:"deptId"`
		RoleID     *uint  `json:"roleId"`
		DeptName   string `json:"deptName"`
		RoleName   string `json:"roleName"`
		Status     int8   `json:"status"`
		CreateTime string `json:"createTime"`
	}

	var list []UserItem
	for _, u := range users {
		item := UserItem{
			ID:         u.ID,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Email:      u.Email,
			Phone:      u.Phone,
			Avatar:     u.Avatar,
			DeptID:     u.DeptID,
			RoleID:     u.RoleID,
			Status:     u.Status,
			CreateTime: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if u.DeptID != nil {
			var dept model.Department
			if err := h.db.Where("id = ? AND tenant_id = ?", *u.DeptID, tid).First(&dept).Error; err == nil {
				item.DeptName = dept.Name
			}
		}
		if u.RoleID != nil {
			var role model.Role
			if err := h.db.Where("id = ? AND tenant_id = ?", *u.RoleID, tid).First(&role).Error; err == nil {
				item.RoleName = role.Name
			}
		}
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"list":  list,
			"total": total,
		},
	})
}

func (h *UserHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Phone    string `json:"phone"`
		DeptID   *uint  `json:"deptId"`
		RoleID   *uint  `json:"roleId"`
		Status   *int8  `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var count int64
	h.db.Model(&model.User{}).Where("username = ? AND tenant_id = ?", req.Username, tid).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "用户名已存在"})
		return
	}

	var tenant model.Tenant
	if err := h.db.First(&tenant, tid).Error; err == nil {
		if tenant.CurrentUsers >= tenant.MaxUsers {
			c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "套餐用户数已满，请升级套餐"})
			return
		}
	}

	hashedPwd, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "密码加密失败"})
		return
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	user := model.User{
		TenantID: tid,
		Username: req.Username,
		Password: hashedPwd,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		DeptID:   req.DeptID,
		RoleID:   req.RoleID,
		Status:   status,
	}

	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建用户失败"})
		return
	}

	h.db.Model(&model.Tenant{}).Where("id = ?", tid).Update("current_users", gorm.Expr("current_users + 1"))

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": user})
}

func (h *UserHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Avatar   string `json:"avatar"`
		DeptID   *uint  `json:"deptId"`
		RoleID   *uint  `json:"roleId"`
		Status   *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	result := h.db.Model(&model.User{}).Where("id = ? AND tenant_id = ?", id, tid).Updates(map[string]interface{}{
		"nickname": req.Nickname,
		"email":    req.Email,
		"phone":    req.Phone,
		"avatar":   req.Avatar,
		"dept_id":  req.DeptID,
		"role_id":  req.RoleID,
		"status":   req.Status,
	})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *UserHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).Delete(&model.User{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	h.db.Model(&model.Tenant{}).Where("id = ?", tid).Update("current_users", gorm.Expr("GREATEST(current_users - 1, 0)"))
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	h.db.Model(&model.User{}).Where("id = ? AND tenant_id = ?", id, tid).Update("status", req.Status)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}
