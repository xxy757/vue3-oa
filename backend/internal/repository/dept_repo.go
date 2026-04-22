package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type DeptRepo struct {
	db *gorm.DB
}

func NewDeptRepo(db *gorm.DB) *DeptRepo {
	return &DeptRepo{db: db}
}

func (r *DeptRepo) ListByTenant(tenantID uint) ([]model.Department, error) {
	var depts []model.Department
	if err := r.db.Where("tenant_id = ?", tenantID).Order("sort ASC").Find(&depts).Error; err != nil {
		return nil, err
	}
	return depts, nil
}

func (r *DeptRepo) GetByID(id uint, tenantID uint) (*model.Department, error) {
	var dept model.Department
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&dept).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

func (r *DeptRepo) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

func (r *DeptRepo) Update(id uint, tenantID uint, updates map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.Department{}).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *DeptRepo) CountChildren(parentID uint, tenantID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Department{}).Where("parent_id = ? AND tenant_id = ?", parentID, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *DeptRepo) Delete(id uint, tenantID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Department{}).Error
}
