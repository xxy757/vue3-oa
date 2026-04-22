package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
	"time"

	"gorm.io/gorm"
)

type AdminService struct {
	tenantRepo *repository.TenantRepo
	userRepo   *repository.UserRepo
	planRepo   *repository.PlanRepo
	db         *gorm.DB
}

func NewAdminService(tenantRepo *repository.TenantRepo, userRepo *repository.UserRepo, planRepo *repository.PlanRepo, db *gorm.DB) *AdminService {
	return &AdminService{
		tenantRepo: tenantRepo,
		userRepo:   userRepo,
		planRepo:   planRepo,
		db:         db,
	}
}

type DashboardData struct {
	TotalTenants      int64          `json:"totalTenants"`
	ActiveTenants     int64          `json:"activeTenants"`
	NewTenantsThisMonth int64        `json:"newTenantsThisMonth"`
	TotalUsers        int64          `json:"totalUsers"`
	MonthlyRevenue    float64        `json:"monthlyRevenue"`
	RevenueGrowth     float64        `json:"revenueGrowth"`
	PlanDistribution  []PlanDistItem `json:"planDistribution"`
	RecentTenants     []RecentTenant `json:"recentTenants"`
}

type PlanDistItem struct {
	Name       string  `json:"name"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
	Color      string  `json:"color"`
}

type RecentTenant struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	CreateTime string `json:"createTime"`
}

// Dashboard 聚合平台管理员仪表盘数据，包括租户统计、营收、套餐分布和最近注册租户
func (s *AdminService) Dashboard() (*DashboardData, error) {
	// 统计租户总数
	totalTenants, _ := s.tenantRepo.CountAll()
	// 统计活跃租户数（状态为 active）
	activeTenants, _ := s.tenantRepo.CountByStatus("active")

	// 计算本月1号零点作为起始时间
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	// 查询本月新增租户数
	newTenantsThisMonth, _ := s.tenantRepo.CountCreatedAfter(monthStart)

	// 全平台用户总数
	var totalUsers int64
	s.db.Model(&model.User{}).Count(&totalUsers)

	// 按套餐 ID 分组统计租户数量
	planCounts, _ := s.tenantRepo.CountByPlan()
	// 查询所有套餐定义，建成 map 方便按 ID 查找套餐名称
	plans, _ := s.planRepo.ListAll()
	planMap := make(map[uint]model.Plan)
	for _, p := range plans {
		planMap[p.ID] = p
	}

	// 套餐编码到前端展示颜色的静态映射
	colorMap := map[string]string{
		"free": "#8c8c8c", "standard": "#1677FF", "professional": "#722ED1", "enterprise": "#FA8C16",
	}
	// 遍历每个套餐的计数结果，组装分布数据
	var planDistribution []PlanDistItem
	for _, pc := range planCounts {
		// 从 planMap 查找套餐名称，找不到则显示"未知"
		p, ok := planMap[pc.PlanID]
		name := "未知"
		color := "#8c8c8c"
		if ok {
			name = p.Name
			// 从 colorMap 查找对应颜色
			if c, exists := colorMap[p.Code]; exists {
				color = c
			}
		}
		// 计算该套餐占总租户的百分比
		var pct float64
		if totalTenants > 0 {
			pct = float64(pc.Count) / float64(totalTenants) * 100
		}
		planDistribution = append(planDistribution, PlanDistItem{
			Name: name, Count: pc.Count, Percentage: pct, Color: color,
		})
	}

	// 本月已付款发票金额总和（nil 表示无截止时间，即到当前）
	monthlyRevenue, _ := s.tenantRepo.SumInvoiceAmount("paid", monthStart, nil)
	// 上月1号零点
	lastMonthStart := monthStart.AddDate(0, -1, 0)
	// 上月已付款发票金额总和
	lastMonthRevenue, _ := s.tenantRepo.SumInvoiceAmount("paid", lastMonthStart, monthStart)
	// 环比增长率：(本月 - 上月) / 上月 * 100；上月为 0 时不计算，保持 0
	var revenueGrowth float64
	if lastMonthRevenue > 0 {
		revenueGrowth = ((monthlyRevenue - lastMonthRevenue) / lastMonthRevenue) * 100
	}

	// 取最新 5 条租户记录
	recentDB, _ := s.tenantRepo.ListRecent(5)
	var recentTenants []RecentTenant
	for _, t := range recentDB {
		// 将 GORM 模型转为 API 返回结构，时间格式化为 ISO 8601
		recentTenants = append(recentTenants, RecentTenant{
			ID: t.ID, Name: t.Name, Status: t.Status,
			CreateTime: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	// 打包所有统计结果返回给前端
	return &DashboardData{
		TotalTenants: totalTenants, ActiveTenants: activeTenants,
		NewTenantsThisMonth: newTenantsThisMonth, TotalUsers: totalUsers,
		MonthlyRevenue: monthlyRevenue, RevenueGrowth: revenueGrowth,
		PlanDistribution: planDistribution, RecentTenants: recentTenants,
	}, nil
}

type TenantItem struct {
	model.Tenant
	Plan model.Plan `json:"plan"`
}

func (s *AdminService) ListTenants() ([]TenantItem, error) {
	tenants, err := s.tenantRepo.ListAll()
	if err != nil {
		return nil, utils.ErrInternal("获取租户列表失败")
	}
	var list []TenantItem
	for _, t := range tenants {
		var plan model.Plan
		if p, err := s.planRepo.GetByID(t.PlanID); err == nil {
			plan = *p
		}
		list = append(list, TenantItem{Tenant: t, Plan: plan})
	}
	return list, nil
}

func (s *AdminService) CreateTenant(name, slug, contactName, contactPhone, contactEmail string, planID uint) (*model.Tenant, error) {
	count, err := s.tenantRepo.CountBySlug(slug)
	if err != nil {
		return nil, utils.ErrInternal("查询企业失败")
	}
	if count > 0 {
		return nil, utils.ErrConflict("企业标识已被占用")
	}
	resolvedPlanID := planID
	if resolvedPlanID == 0 {
		if freePlan, err := s.planRepo.GetByCode("free"); err == nil {
			resolvedPlanID = freePlan.ID
		}
	}
	plan, err := s.planRepo.GetByID(resolvedPlanID)
	if err != nil {
		return nil, utils.ErrNotFound("套餐不存在")
	}
	trialEnds := time.Now().Add(14 * 24 * time.Hour)
	tenant := &model.Tenant{
		Name: name, Slug: slug,
		ContactName: contactName, ContactPhone: contactPhone, ContactEmail: contactEmail,
		PlanID: resolvedPlanID, MaxUsers: plan.MaxUsers, Status: "trial", TrialEndsAt: &trialEnds,
	}
	if err := s.tenantRepo.Create(tenant); err != nil {
		return nil, utils.ErrInternal("创建租户失败")
	}
	return tenant, nil
}

func (s *AdminService) UpdateTenant(id string, updates map[string]interface{}) error {
	return s.tenantRepo.UpdateByStringID(id, updates)
}

func (s *AdminService) ActivateTenant(id string) error {
	return s.tenantRepo.UpdateByStringID(id, map[string]interface{}{"status": "active"})
}

func (s *AdminService) SuspendTenant(id string) error {
	return s.tenantRepo.UpdateByStringID(id, map[string]interface{}{"status": "suspended"})
}

func (s *AdminService) ListPlans() ([]model.Plan, error) {
	return s.planRepo.ListAll()
}

func (s *AdminService) CreatePlan(name, code string, price float64, minUsers, maxUsers int, features model.FeatureMap) (*model.Plan, error) {
	plan := &model.Plan{
		Name: name, Code: code, Price: price,
		MinUsers: minUsers, MaxUsers: maxUsers,
		Features: features, IsActive: 1,
	}
	if err := s.planRepo.Create(plan); err != nil {
		return nil, utils.ErrInternal("创建套餐失败")
	}
	return plan, nil
}

func (s *AdminService) UpdatePlan(id string, updates map[string]interface{}) error {
	return s.planRepo.Update(id, updates)
}
