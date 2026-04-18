package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	var totalTenants int64
	h.db.Model(&model.Tenant{}).Count(&totalTenants)

	var activeTenants int64
	h.db.Model(&model.Tenant{}).Where("status = ?", "active").Count(&activeTenants)

	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	var newTenantsThisMonth int64
	h.db.Model(&model.Tenant{}).Where("created_at >= ?", monthStart).Count(&newTenantsThisMonth)

	var totalUsers int64
	h.db.Model(&model.User{}).Count(&totalUsers)

	type PlanCount struct {
		PlanID uint `json:"planId"`
		Count  int64 `json:"count"`
	}
	var planCounts []PlanCount
	h.db.Model(&model.Tenant{}).Select("plan_id, count(*) as count").Group("plan_id").Find(&planCounts)

	var monthlyRevenue float64
	h.db.Model(&model.Invoice{}).
		Where("status = ? AND created_at >= ?", "paid", monthStart).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthlyRevenue)

	var lastMonthRevenue float64
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	h.db.Model(&model.Invoice{}).
		Where("status = ? AND created_at >= ? AND created_at < ?", "paid", lastMonthStart, monthStart).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&lastMonthRevenue)

	var revenueGrowth float64
	if lastMonthRevenue > 0 {
		revenueGrowth = ((monthlyRevenue - lastMonthRevenue) / lastMonthRevenue) * 100
	}

	var plans []model.Plan
	h.db.Find(&plans)
	planMap := make(map[uint]model.Plan)
	for _, p := range plans {
		planMap[p.ID] = p
	}

	type PlanDistItem struct {
		Name       string  `json:"name"`
		Count      int64   `json:"count"`
		Percentage float64 `json:"percentage"`
		Color      string  `json:"color"`
	}
	var planDistribution []PlanDistItem
	colorMap := map[string]string{
		"free": "#8c8c8c", "standard": "#1677FF", "professional": "#722ED1", "enterprise": "#FA8C16",
	}
	for _, pc := range planCounts {
		p, ok := planMap[pc.PlanID]
		name := "未知"
		color := "#8c8c8c"
		if ok {
			name = p.Name
			if c, exists := colorMap[p.Code]; exists {
				color = c
			}
		}
		var pct float64
		if totalTenants > 0 {
			pct = float64(pc.Count) / float64(totalTenants) * 100
		}
		planDistribution = append(planDistribution, PlanDistItem{
			Name: name, Count: pc.Count, Percentage: pct, Color: color,
		})
	}

	type RecentTenant struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		Status     string `json:"status"`
		CreateTime string `json:"createTime"`
	}
	var recentDB []model.Tenant
	h.db.Order("created_at DESC").Limit(5).Find(&recentDB)
	var recentTenants []RecentTenant
	for _, t := range recentDB {
		recentTenants = append(recentTenants, RecentTenant{
			ID: t.ID, Name: t.Name, Status: t.Status,
			CreateTime: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
		"totalTenants":      totalTenants,
		"activeTenants":     activeTenants,
		"newTenantsThisMonth": newTenantsThisMonth,
		"totalUsers":        totalUsers,
		"monthlyRevenue":    monthlyRevenue,
		"revenueGrowth":     revenueGrowth,
		"planDistribution":  planDistribution,
		"recentTenants":     recentTenants,
	}})
}

func (h *AdminHandler) ListTenants(c *gin.Context) {
	var tenants []model.Tenant
	h.db.Order("created_at DESC").Find(&tenants)

	type TenantItem struct {
		model.Tenant
		Plan model.Plan `json:"plan"`
	}
	var list []TenantItem
	for _, t := range tenants {
		var plan model.Plan
		h.db.First(&plan, t.PlanID)
		list = append(list, TenantItem{Tenant: t, Plan: plan})
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
		Name: req.Name, Slug: req.Slug,
		ContactName: req.ContactName, ContactPhone: req.ContactPhone, ContactEmail: req.ContactEmail,
		PlanID: planID, MaxUsers: plan.MaxUsers, Status: "trial", TrialEndsAt: &trialEnds,
	}
	if err := h.db.Create(&tenant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建租户失败"})
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
	h.db.Model(&model.Tenant{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *AdminHandler) ActivateTenant(c *gin.Context) {
	id := c.Param("id")
	h.db.Model(&model.Tenant{}).Where("id = ?", id).Update("status", "active")
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "启用成功"})
}

func (h *AdminHandler) SuspendTenant(c *gin.Context) {
	id := c.Param("id")
	h.db.Model(&model.Tenant{}).Where("id = ?", id).Update("status", "suspended")
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "暂停成功"})
}

func (h *AdminHandler) ListPlans(c *gin.Context) {
	var plans []model.Plan
	h.db.Order("price ASC").Find(&plans)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": plans})
}

func (h *AdminHandler) CreatePlan(c *gin.Context) {
	var req struct {
		Name     string          `json:"name" binding:"required"`
		Code     string          `json:"code" binding:"required"`
		Price    float64         `json:"price"`
		MinUsers int             `json:"minUsers"`
		MaxUsers int             `json:"maxUsers"`
		Features model.FeatureMap `json:"features"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	plan := model.Plan{
		Name: req.Name, Code: req.Code, Price: req.Price,
		MinUsers: req.MinUsers, MaxUsers: req.MaxUsers,
		Features: req.Features, IsActive: 1,
	}
	if err := h.db.Create(&plan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建套餐失败"})
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
	h.db.Model(&model.Plan{}).Where("id = ?", id).Updates(updates)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}
