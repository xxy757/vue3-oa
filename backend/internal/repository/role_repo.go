package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) ListByTenant(tenantID uint) ([]model.Role, error) {
	var roles []model.Role
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepo) GetByID(id uint, tenantID uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) GetByCode(code string, tenantID uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.Where("code = ? AND tenant_id = ?", code, tenantID).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepo) CountByCode(code string, tenantID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Role{}).Where("code = ? AND tenant_id = ?", code, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RoleRepo) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) Update(id uint, tenantID uint, updates map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.Role{}).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *RoleRepo) Delete(role *model.Role) error {
	return r.db.Delete(role).Error
}
