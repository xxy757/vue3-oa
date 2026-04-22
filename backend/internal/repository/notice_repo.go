package repository

import (
	"oa-saas/internal/model"

	"gorm.io/gorm"
)

type NoticeRepo struct {
	db *gorm.DB
}

func NewNoticeRepo(db *gorm.DB) *NoticeRepo {
	return &NoticeRepo{db: db}
}

func (r *NoticeRepo) ListByTenant(tenantID uint, noticeType int, keyword string, page, pageSize int) ([]model.Notice, int64, error) {
	var notices []model.Notice
	var total int64
	query := r.db.Model(&model.Notice{}).Where("tenant_id = ?", tenantID)
	if noticeType > 0 {
		query = query.Where("type = ?", noticeType)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("is_top DESC, created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&notices).Error; err != nil {
		return nil, 0, err
	}
	return notices, total, nil
}

func (r *NoticeRepo) GetByID(id uint, tenantID uint) (*model.Notice, error) {
	var notice model.Notice
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&notice).Error; err != nil {
		return nil, err
	}
	return &notice, nil
}

func (r *NoticeRepo) Create(notice *model.Notice) error {
	return r.db.Create(notice).Error
}

func (r *NoticeRepo) GetOrCreateRead(noticeID, userID, tenantID uint, read *model.NoticeRead) error {
	return r.db.Where("notice_id = ? AND user_id = ? AND tenant_id = ?", noticeID, userID, tenantID).
		FirstOrCreate(read, model.NoticeRead{NoticeID: noticeID, UserID: userID, TenantID: tenantID}).Error
}

func (r *NoticeRepo) HasRead(noticeID, userID, tenantID uint) (bool, error) {
	var read model.NoticeRead
	result := r.db.Where("notice_id = ? AND user_id = ? AND tenant_id = ?", noticeID, userID, tenantID).First(&read)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}
	return true, nil
}

func (r *NoticeRepo) CountUnread(userID, tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notice{}).
		Where("tenant_id = ? AND id NOT IN (SELECT notice_id FROM notice_reads WHERE user_id = ? AND tenant_id = ?)", tenantID, userID, tenantID).
		Count(&count).Error
	return count, err
}
