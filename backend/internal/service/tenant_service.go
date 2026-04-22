package service

import (
	"fmt"
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
	"time"

	"gorm.io/gorm"
)

type TenantService struct {
	tenantRepo *repository.TenantRepo
	roleRepo   *repository.RoleRepo
	userRepo   *repository.UserRepo
	planRepo   *repository.PlanRepo
	db         *gorm.DB
}

func NewTenantService(tenantRepo *repository.TenantRepo, roleRepo *repository.RoleRepo, userRepo *repository.UserRepo, planRepo *repository.PlanRepo, db *gorm.DB) *TenantService {
	return &TenantService{
		tenantRepo: tenantRepo,
		roleRepo:   roleRepo,
		userRepo:   userRepo,
		planRepo:   planRepo,
		db:         db,
	}
}

type RegisterResult struct {
	TenantID    uint      `json:"tenantId"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	TrialEndsAt *time.Time `json:"trialEndsAt"`
	AdminUser   struct {
		ID           uint   `json:"id"`
		Username     string `json:"username"`
		TempPassword string `json:"tempPassword"`
	} `json:"adminUser"`
}

func (s *TenantService) Register(name, slug, contactName, contactPhone, contactEmail string, planID uint) (*RegisterResult, error) {
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

	var result RegisterResult
	err = s.db.Transaction(func(tx *gorm.DB) error {
		txTenantRepo := repository.NewTenantRepo(tx)
		txRoleRepo := repository.NewRoleRepo(tx)
		txUserRepo := repository.NewUserRepo(tx)

		if err := txTenantRepo.Create(tenant); err != nil {
			return err
		}
		adminRole := &model.Role{
			TenantID: tenant.ID, Name: "管理员", Code: "admin",
			Description: "系统管理员", Permissions: model.StringArray{"*"}, Status: 1,
		}
		if err := txRoleRepo.Create(adminRole); err != nil {
			return err
		}
		userRole := &model.Role{
			TenantID: tenant.ID, Name: "普通员工", Code: "employee",
			Description: "普通员工角色", Permissions: model.StringArray{"approval:apply", "notice:view", "schedule:view"}, Status: 1,
		}
		if err := txRoleRepo.Create(userRole); err != nil {
			return err
		}
		tempPwd := "Abc123456"
		hashedPwd, _ := utils.HashPassword(tempPwd)
		admin := &model.User{
			TenantID: tenant.ID, Username: "admin", Password: hashedPwd,
			Nickname: contactName, Email: contactEmail, Phone: contactPhone,
			RoleID: &adminRole.ID, Status: 1,
		}
		if err := txUserRepo.Create(admin); err != nil {
			return err
		}
		if err := txTenantRepo.IncrementCurrentUsers(tenant.ID); err != nil {
			return err
		}
		result = RegisterResult{
			TenantID:    tenant.ID,
			Name:        tenant.Name,
			Slug:        tenant.Slug,
			TrialEndsAt: tenant.TrialEndsAt,
		}
		result.AdminUser.ID = admin.ID
		result.AdminUser.Username = admin.Username
		result.AdminUser.TempPassword = tempPwd
		return nil
	})
	if err != nil {
		return nil, utils.ErrInternal("注册失败")
	}
	return &result, nil
}

type TenantInfoResult struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Slug         string     `json:"slug"`
	Logo         string     `json:"logo"`
	ContactName  string     `json:"contactName"`
	ContactPhone string     `json:"contactPhone"`
	ContactEmail string     `json:"contactEmail"`
	CurrentUsers int        `json:"currentUsers"`
	MaxUsers     int        `json:"maxUsers"`
	Status       string     `json:"status"`
	TrialEndsAt  *time.Time `json:"trialEndsAt"`
	PlanExpireAt *time.Time `json:"planExpireAt"`
	Plan         struct {
		ID       uint             `json:"id"`
		Name     string           `json:"name"`
		Code     string           `json:"code"`
		Price    float64          `json:"price"`
		Features model.FeatureMap `json:"features"`
		MaxUsers int              `json:"maxUsers"`
	} `json:"plan"`
}

func (s *TenantService) GetInfo(tenantID uint) (*TenantInfoResult, error) {
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("租户不存在")
	}
	plan, _ := s.planRepo.GetByID(tenant.PlanID)
	result := &TenantInfoResult{
		ID: tenant.ID, Name: tenant.Name, Slug: tenant.Slug,
		Logo: tenant.Logo, ContactName: tenant.ContactName,
		ContactPhone: tenant.ContactPhone, ContactEmail: tenant.ContactEmail,
		CurrentUsers: tenant.CurrentUsers, MaxUsers: tenant.MaxUsers,
		Status: tenant.Status, TrialEndsAt: tenant.TrialEndsAt, PlanExpireAt: tenant.PlanExpireAt,
	}
	if plan != nil {
		result.Plan.ID = plan.ID
		result.Plan.Name = plan.Name
		result.Plan.Code = plan.Code
		result.Plan.Price = plan.Price
		result.Plan.Features = plan.Features
		result.Plan.MaxUsers = plan.MaxUsers
	}
	return result, nil
}

func (s *TenantService) UpdateInfo(tenantID uint, updates map[string]interface{}) error {
	return s.tenantRepo.Update(tenantID, updates)
}

func (s *TenantService) ListPlans() ([]model.Plan, error) {
	plans, err := s.planRepo.ListActive()
	if err != nil {
		return nil, utils.ErrInternal("获取套餐列表失败")
	}
	return plans, nil
}

type UpgradeResult struct {
	PlanID    uint      `json:"planId"`
	PlanName  string    `json:"planName"`
	ExpireAt  time.Time `json:"expireAt"`
	InvoiceNo string    `json:"invoiceNo"`
}

func (s *TenantService) UpgradePlan(tenantID uint, planID uint) (*UpgradeResult, error) {
	plan, err := s.planRepo.GetByID(planID)
	if err != nil {
		return nil, utils.ErrNotFound("套餐不存在")
	}
	tenant, _ := s.tenantRepo.GetByID(tenantID)
	now := time.Now()
	expireAt := now.Add(30 * 24 * time.Hour)
	updates := map[string]interface{}{
		"plan_id": planID, "max_users": plan.MaxUsers,
		"plan_start_at": now, "plan_expire_at": expireAt, "status": "active",
	}
	if err := s.tenantRepo.Update(tenantID, updates); err != nil {
		return nil, utils.ErrInternal("升级失败")
	}
	invoiceNo := fmt.Sprintf("INV-%d-%d", tenantID, now.Unix())
	var currentUsers int
	if tenant != nil {
		currentUsers = tenant.CurrentUsers
	}
	invoice := &model.Invoice{
		TenantID: tenantID, PlanID: planID, InvoiceNo: invoiceNo,
		PeriodStart: now, PeriodEnd: expireAt, UserCount: currentUsers,
		Amount: plan.Price * float64(currentUsers), Status: "paid", PaidAt: &now,
	}
	s.tenantRepo.CreateInvoice(invoice)
	return &UpgradeResult{
		PlanID: planID, PlanName: plan.Name, ExpireAt: expireAt, InvoiceNo: invoiceNo,
	}, nil
}

func (s *TenantService) ListInvoices(tenantID uint) ([]model.Invoice, error) {
	invoices, err := s.tenantRepo.ListInvoices(tenantID)
	if err != nil {
		return nil, utils.ErrInternal("获取账单失败")
	}
	return invoices, nil
}
