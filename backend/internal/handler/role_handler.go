package handler

import (
	"net/http"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

func (h *RoleHandler) List(c *gin.Context) {
	tid := getTenantID(c)
	roles, err := h.roleService.List(tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": roles})
}

func (h *RoleHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
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
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	role, err := h.roleService.Create(tid, req.Name, req.Code, req.Description, req.Permissions, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": role})
}

func (h *RoleHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
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
		"name": req.Name, "description": req.Description, "permissions": req.Permissions,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if err := h.roleService.Update(uint(id), tid, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *RoleHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.roleService.Delete(uint(id), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
