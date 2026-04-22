package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/service"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminService *service.AdminService
}

func NewAdminHandler(adminService *service.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	result, err := h.adminService.Dashboard()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *AdminHandler) ListTenants(c *gin.Context) {
	list, err := h.adminService.ListTenants()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": list})
}

func (h *AdminHandler) CreateTenant(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		Slug         string `json:"slug" binding:"required"`
		ContactName  string `json:"contactName" binding:"required"`
		ContactPhone string `json:"contactPhone" binding:"required"`
		ContactEmail string `json:"contactEmail" binding:"required"`
		PlanID       uint   `json:"planId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	tenant, err := h.adminService.CreateTenant(req.Name, req.Slug, req.ContactName, req.ContactPhone, req.ContactEmail, req.PlanID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": tenant})
}

func (h *AdminHandler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name         string `json:"name"`
		ContactName  string `json:"contactName"`
		ContactPhone string `json:"contactPhone"`
		ContactEmail string `json:"contactEmail"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ContactName != "" {
		updates["contact_name"] = req.ContactName
	}
	if req.ContactPhone != "" {
		updates["contact_phone"] = req.ContactPhone
	}
	if req.ContactEmail != "" {
		updates["contact_email"] = req.ContactEmail
	}
	if err := h.adminService.UpdateTenant(id, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *AdminHandler) ActivateTenant(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminService.ActivateTenant(id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "启用成功"})
}

func (h *AdminHandler) SuspendTenant(c *gin.Context) {
	id := c.Param("id")
	if err := h.adminService.SuspendTenant(id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "暂停成功"})
}

func (h *AdminHandler) ListPlans(c *gin.Context) {
	plans, err := h.adminService.ListPlans()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plans})
}

func (h *AdminHandler) CreatePlan(c *gin.Context) {
	var req struct {
		Name     string           `json:"name" binding:"required"`
		Code     string           `json:"code" binding:"required"`
		Price    float64          `json:"price"`
		MinUsers int              `json:"minUsers"`
		MaxUsers int              `json:"maxUsers"`
		Features model.FeatureMap `json:"features"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	plan, err := h.adminService.CreatePlan(req.Name, req.Code, req.Price, req.MinUsers, req.MaxUsers, req.Features)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plan})
}

func (h *AdminHandler) UpdatePlan(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name     *string          `json:"name"`
		Price    *float64         `json:"price"`
		IsActive *int8            `json:"isActive"`
		MinUsers *int             `json:"minUsers"`
		MaxUsers *int             `json:"maxUsers"`
		Features model.FeatureMap `json:"features"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	if req.MinUsers != nil {
		updates["min_users"] = *req.MinUsers
	}
	if req.MaxUsers != nil {
		updates["max_users"] = *req.MaxUsers
	}
	if req.Features != nil {
		updates["features"] = req.Features
	}
	if err := h.adminService.UpdatePlan(id, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}
