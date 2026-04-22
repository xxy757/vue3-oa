package repository

import (
	"oa-saas/internal/model"
	"time"

	"gorm.io/gorm"
)

type ScheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) *ScheduleRepo {
	return &ScheduleRepo{db: db}
}

func (r *ScheduleRepo) GetByID(id uint, tenantID uint) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *ScheduleRepo) Create(schedule *model.Schedule) error {
	return r.db.Create(schedule).Error
}

func (r *ScheduleRepo) Update(id uint, tenantID uint, updates map[string]interface{}) (int64, error) {
	result := r.db.Model(&model.Schedule{}).Where("id = ? AND tenant_id = ?", id, tenantID).Updates(updates)
	return result.RowsAffected, result.Error
}

func (r *ScheduleRepo) Delete(id uint, tenantID uint) error {
	r.db.Where("schedule_id = ? AND tenant_id = ?", id, tenantID).Delete(&model.ScheduleParticipant{})
	return r.db.Where("id = ? AND tenant_id = ?", id, tenantID).Delete(&model.Schedule{}).Error
}

func (r *ScheduleRepo) ListByUserAndDateRange(userID, tenantID uint, startDate, endDate *time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	query := r.db.Model(&model.Schedule{}).
		Joins("LEFT JOIN schedule_participants ON schedule_participants.schedule_id = schedules.id").
		Where("(schedules.creator_id = ? OR schedule_participants.user_id = ?) AND schedules.tenant_id = ?", userID, userID, tenantID)
	if startDate != nil && endDate != nil {
		query = query.Where("schedules.start_time <= ? AND schedules.end_time >= ?", *endDate, *startDate)
	}
	if err := query.Distinct("schedules.id").Order("schedules.start_time ASC").Find(&schedules).Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r *ScheduleRepo) ListByUserAndWeek(userID, tenantID uint, startOfWeek, endOfWeek time.Time) ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Model(&model.Schedule{}).
		Joins("LEFT JOIN schedule_participants ON schedule_participants.schedule_id = schedules.id").
		Where("(schedules.creator_id = ? OR schedule_participants.user_id = ?) AND schedules.tenant_id = ? AND schedules.start_time >= ? AND schedules.start_time < ?", userID, userID, tenantID, startOfWeek, endOfWeek).
		Group("schedules.id").
		Order("schedules.start_time ASC").
		Find(&schedules).Error
	return schedules, err
}

func (r *ScheduleRepo) ListParticipants(scheduleID, tenantID uint) ([]model.ScheduleParticipant, error) {
	var participants []model.ScheduleParticipant
	if err := r.db.Where("schedule_id = ? AND tenant_id = ?", scheduleID, tenantID).Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ScheduleRepo) CreateParticipant(participant *model.ScheduleParticipant) error {
	return r.db.Create(participant).Error
}
