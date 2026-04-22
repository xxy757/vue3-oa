package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type TenantRepo struct {
	db *gorm.DB
}

func NewTenantRepo(db *gorm.DB) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) GetByID(id uint) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.First(&tenant, id).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) GetBySlug(slug string) (*model.Tenant, error) {
	var tenant model.Tenant
	if err := r.db.Where("slug = ?", slug).First(&tenant).Error; err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (r *TenantRepo) CountBySlug(slug string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Tenant{}).Where("slug = ?", slug).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TenantRepo) Create(tenant *model.Tenant) error {
	return r.db.Create(tenant).Error
}

func (r *TenantRepo) Update(id uint, updates map[string]interface{}) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TenantRepo) IncrementCurrentUsers(id uint) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Update("current_users", gorm.Expr("current_users + 1")).Error
}

func (r *TenantRepo) DecrementCurrentUsers(id uint) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Update("current_users", gorm.Expr("GREATEST(current_users - 1, 0)")).Error
}

func (r *TenantRepo) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Tenant{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TenantRepo) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Tenant{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TenantRepo) CountCreatedAfter(t interface{}) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Tenant{}).Where("created_at >= ?", t).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *TenantRepo) ListAll() ([]model.Tenant, error) {
	var tenants []model.Tenant
	if err := r.db.Order("created_at DESC").Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepo) ListRecent(limit int) ([]model.Tenant, error) {
	var tenants []model.Tenant
	if err := r.db.Order("created_at DESC").Limit(limit).Find(&tenants).Error; err != nil {
		return nil, err
	}
	return tenants, nil
}

func (r *TenantRepo) CountByPlan() ([]struct {
	PlanID uint
	Count  int64
}, error) {
	var result []struct {
		PlanID uint
		Count  int64
	}
	if err := r.db.Model(&model.Tenant{}).Select("plan_id, count(*) as count").Group("plan_id").Find(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (r *TenantRepo) CreateInvoice(invoice *model.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *TenantRepo) ListInvoices(tenantID uint) ([]model.Invoice, error) {
	var invoices []model.Invoice
	if err := r.db.Where("tenant_id = ?", tenantID).Order("created_at DESC").Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *TenantRepo) UpdateByStringID(id string, updates map[string]interface{}) error {
	return r.db.Model(&model.Tenant{}).Where("id = ?", id).Updates(updates).Error
}

func (r *TenantRepo) SumInvoiceAmount(status string, after interface{}, before interface{}) (float64, error) {
	var amount float64
	q := r.db.Model(&model.Invoice{}).Where("status = ?", status)
	if after != nil {
		q = q.Where("created_at >= ?", after)
	}
	if before != nil {
		q = q.Where("created_at < ?", before)
	}
	if err := q.Select("COALESCE(SUM(amount), 0)").Scan(&amount).Error; err != nil {
		return 0, err
	}
	return amount, nil
}
