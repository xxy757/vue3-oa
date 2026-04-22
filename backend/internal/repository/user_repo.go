package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(id uint, tenantID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByUsername(username string, tenantID uint) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ? AND tenant_id = ?", username, tenantID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) ListByTenant(tenantID uint, keyword string, page, pageSize int) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := r.db.Model(&model.User{}).Where("tenant_id = ?", tenantID)
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *UserRepo) CountByUsername(username string, tenantID uint) (int64, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("username = ? AND tenant_id = ?", username, tenantID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) Update(id uint, tenantID uint, updates map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.User{}).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *UserRepo) Delete(id uint, tenantID uint) error {
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.User{}).Error
}

func (r *UserRepo) UpdatePassword(id uint, tenantID uint, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ? AND tenant_id = ?", id, tenantID).Update("password", hashedPassword).Error
}

func (r *UserRepo) UpdateStatus(id uint, tenantID uint, status int8) error {
	return r.db.Model(&model.User{}).Where("id = ? AND tenant_id = ?", id, tenantID).Update("status", status).Error
}
