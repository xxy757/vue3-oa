package handler

import (
	"fmt"
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TenantHandler struct {
	db *gorm.DB
}

func NewTenantHandler(db *gorm.DB) *TenantHandler {
	return &TenantHandler{db: db}
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

	var count int64
	h.db.Model(&model.Tenant{}).Where("slug = ?", req.Slug).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "企业标识已被占用"})
		return
	}

	planID := req.PlanID
	if planID == 0 {
		var freePlan model.Plan
		if err := h.db.Where("code = ?", "free").First(&freePlan).Error; err == nil {
			planID = freePlan.ID
		}
	}

	var plan model.Plan
	if err := h.db.First(&plan, planID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "套餐不存在"})
		return
	}

	trialEnds := time.Now().Add(14 * 24 * time.Hour)
	tenant := model.Tenant{
		Name:         req.Name,
		Slug:         req.Slug,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		ContactEmail: req.ContactEmail,
		PlanID:       planID,
		MaxUsers:     plan.MaxUsers,
		Status:       "trial",
		TrialEndsAt:  &trialEnds,
	}
	if err := h.db.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "注册失败"})
		return
	}

	adminRole := model.Role{
		TenantID:    tenant.ID,
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员",
		Permissions: model.StringArray{"*"},
		Status:      1,
	}
	h.db.Create(&adminRole)

	userRole := model.Role{
		TenantID:    tenant.ID,
		Name:        "普通员工",
		Code:        "employee",
		Description: "普通员工角色",
		Permissions: model.StringArray{"approval:apply", "notice:view", "schedule:view"},
		Status:      1,
	}
	h.db.Create(&userRole)

	tempPwd := "Abc123456"
	hashedPwd, _ := utils.HashPassword(tempPwd)
	admin := model.User{
		TenantID: tenant.ID,
		Username: "admin",
		Password: hashedPwd,
		Nickname: req.ContactName,
		Email:    req.ContactEmail,
		Phone:    req.ContactPhone,
		RoleID:   &adminRole.ID,
		Status:   1,
	}
	h.db.Create(&admin)
	h.db.Model(&model.Tenant{}).Where("id = ?", tenant.ID).Update("current_users", 1)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"tenantId":    tenant.ID,
			"name":        tenant.Name,
			"slug":        tenant.Slug,
			"trialEndsAt": tenant.TrialEndsAt,
			"adminUser": gin.H{
				"id":           admin.ID,
				"username":     admin.Username,
				"tempPassword": tempPwd,
			},
		},
	})
}

func (h *TenantHandler) GetInfo(c *gin.Context) {
	tid := getTenantID(c)
	var tenant model.Tenant
	if err := h.db.First(&tenant, tid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "租户不存在"})
		return
	}
	var plan model.Plan
	h.db.First(&plan, tenant.PlanID)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":           tenant.ID,
			"name":         tenant.Name,
			"slug":         tenant.Slug,
			"logo":         tenant.Logo,
			"contactName":  tenant.ContactName,
			"contactPhone": tenant.ContactPhone,
			"contactEmail": tenant.ContactEmail,
			"currentUsers": tenant.CurrentUsers,
			"maxUsers":     tenant.MaxUsers,
			"status":       tenant.Status,
			"trialEndsAt":  tenant.TrialEndsAt,
			"planExpireAt": tenant.PlanExpireAt,
			"plan": gin.H{
				"id":       plan.ID,
				"name":     plan.Name,
				"code":     plan.Code,
				"price":    plan.Price,
				"features": plan.Features,
				"maxUsers": plan.MaxUsers,
			},
		},
	})
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

	h.db.Model(&model.Tenant{}).Where("id = ?", tid).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *TenantHandler) ListPlans(c *gin.Context) {
	var plans []model.Plan
	h.db.Where("is_active = 1").Find(&plans)
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

	var plan model.Plan
	if err := h.db.First(&plan, req.PlanID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "套餐不存在"})
		return
	}

	var tenant model.Tenant
	h.db.First(&tenant, tid)

	now := time.Now()
	expireAt := now.Add(30 * 24 * time.Hour)

	h.db.Model(&model.Tenant{}).Where("id = ?", tid).Updates(map[string]interface{}{
		"plan_id":        req.PlanID,
		"max_users":      plan.MaxUsers,
		"plan_start_at":  now,
		"plan_expire_at": expireAt,
		"status":         "active",
	})

	invoiceNo := fmt.Sprintf("INV-%d-%d", tid, now.Unix())
	invoice := model.Invoice{
		TenantID:    tid,
		PlanID:      req.PlanID,
		InvoiceNo:   invoiceNo,
		PeriodStart: now,
		PeriodEnd:   expireAt,
		UserCount:   tenant.CurrentUsers,
		Amount:      plan.Price * float64(tenant.CurrentUsers),
		Status:      "paid",
		PaidAt:      &now,
	}
	h.db.Create(&invoice)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
		"planId":    req.PlanID,
		"planName":  plan.Name,
		"expireAt":  expireAt,
		"invoiceNo": invoiceNo,
	}})
}

func (h *TenantHandler) ListInvoices(c *gin.Context) {
	tid := getTenantID(c)
	var invoices []model.Invoice
	h.db.Where("tenant_id = ?", tid).Order("created_at DESC").Find(&invoices)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": invoices})
}
