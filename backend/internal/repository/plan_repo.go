package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type PlanRepo struct {
	db *gorm.DB
}

func NewPlanRepo(db *gorm.DB) *PlanRepo {
	return &PlanRepo{db: db}
}

func (r *PlanRepo) GetByID(id uint) (*model.Plan, error) {
	var plan model.Plan
	if err := r.db.First(&plan, id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepo) GetByCode(code string) (*model.Plan, error) {
	var plan model.Plan
	if err := r.db.Where("code = ?", code).First(&plan).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepo) ListActive() ([]model.Plan, error) {
	var plans []model.Plan
	if err := r.db.Where("is_active = 1").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepo) ListAll() ([]model.Plan, error) {
	var plans []model.Plan
	if err := r.db.Order("price ASC").Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepo) Create(plan *model.Plan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepo) Update(id string, updates map[string]interface{}) error {
	return r.db.Model(&model.Plan{}).Where("id = ?", id).Updates(updates).Error
}
