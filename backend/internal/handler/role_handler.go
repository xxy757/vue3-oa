package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleHandler struct {
	db *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{db: db}
}

func (h *RoleHandler) List(c *gin.Context) {
	var roles []model.Role
	if err := h.db.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取角色列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": roles})
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code" binding:"required"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions" binding:"required"`
		Status      *int8    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var count int64
	h.db.Model(&model.Role{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "角色编码已存在"})
		return
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	role := model.Role{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Permissions: req.Permissions,
		Status:      status,
	}
	if err := h.db.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建角色失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": role})
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Code        string   `json:"code"`
		Description string   `json:"description"`
		Permissions []string `json:"permissions"`
		Status      *int8    `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"permissions": req.Permissions,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	result := h.db.Model(&model.Role{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role model.Role
	if err := h.db.First(&role, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "角色不存在"})
		return
	}
	if role.Code == "admin" {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无法删除管理员角色"})
		return
	}
	if err := h.db.Delete(&model.Role{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
