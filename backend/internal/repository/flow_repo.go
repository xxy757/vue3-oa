package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type FlowRepo struct {
	db *gorm.DB
}

func NewFlowRepo(db *gorm.DB) *FlowRepo {
	return &FlowRepo{db: db}
}

func (r *FlowRepo) ListByTenant(tenantID uint) ([]model.ApprovalFlow, error) {
	var flows []model.ApprovalFlow
	if err := r.db.Where("tenant_id = ?", tenantID).Find(&flows).Error; err != nil {
		return nil, err
	}
	return flows, nil
}

func (r *FlowRepo) GetByID(id uint, tenantID uint) (*model.ApprovalFlow, error) {
	var flow model.ApprovalFlow
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&flow).Error; err != nil {
		return nil, err
	}
	return &flow, nil
}

func (r *FlowRepo) GetByCode(code string, tenantID uint) (*model.ApprovalFlow, error) {
	var flow model.ApprovalFlow
	if err := r.db.Where("code = ? AND tenant_id = ?", code, tenantID).First(&flow).Error; err != nil {
		return nil, err
	}
	return &flow, nil
}

func (r *FlowRepo) CountByCode(code string, tenantID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.ApprovalFlow{}).Where("code = ? AND tenant_id = ?", code, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *FlowRepo) Create(flow *model.ApprovalFlow) error {
	return r.db.Create(flow).Error
}

func (r *FlowRepo) Update(id uint, tenantID uint, updates map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.ApprovalFlow{}).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *FlowRepo) Delete(id uint, tenantID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.ApprovalFlow{}).Error
}
