package handler

import (
	"net/http"
	"oa-saas/internal/service"

	"github.com/gin-gonic/gin"
)

type TenantHandler struct {
	tenantService *service.TenantService
}

func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{tenantService: tenantService}
}

func (h *TenantHandler) Register(c *gin.Context) {
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
	result, err := h.tenantService.Register(req.Name, req.Slug, req.ContactName, req.ContactPhone, req.ContactEmail, req.PlanID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *TenantHandler) GetInfo(c *gin.Context) {
	tid := getTenantID(c)
	result, err := h.tenantService.GetInfo(tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *TenantHandler) UpdateInfo(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		Name         string `json:"name"`
		Logo         string `json:"logo"`
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
	if req.Logo != "" {
		updates["logo"] = req.Logo
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
	if err := h.tenantService.UpdateInfo(tid, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *TenantHandler) ListPlans(c *gin.Context) {
	plans, err := h.tenantService.ListPlans()
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plans})
}

func (h *TenantHandler) UpgradePlan(c *gin.Context) {
	tid := getTenantID(c)
	var req struct {
		PlanID uint `json:"planId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	result, err := h.tenantService.UpgradePlan(tid, req.PlanID)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *TenantHandler) ListInvoices(c *gin.Context) {
	tid := getTenantID(c)
	invoices, err := h.tenantService.ListInvoices(tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": invoices})
}
