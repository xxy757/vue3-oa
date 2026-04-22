package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ScheduleService struct {
	scheduleRepo *repository.ScheduleRepo
	db           *gorm.DB
}

func NewScheduleService(scheduleRepo *repository.ScheduleRepo, db *gorm.DB) *ScheduleService {
	return &ScheduleService{
		scheduleRepo: scheduleRepo,
		db:           db,
	}
}

func (s *ScheduleService) List(userID, tenantID uint, startDate, endDate string) ([]model.Schedule, error) {
	var start, end *time.Time
	if startDate != "" && endDate != "" {
		sParsed, _ := time.Parse("2006-01-02", startDate)
		eParsed, _ := time.Parse("2006-01-02", endDate)
		start = &sParsed
		end = &eParsed
	}
	schedules, err := s.scheduleRepo.ListByUserAndDateRange(userID, tenantID, start, end)
	if err != nil {
		return nil, utils.ErrInternal("获取日程列表失败")
	}
	return schedules, nil
}

type ScheduleDetail struct {
	Schedule     model.Schedule             `json:"schedule"`
	Participants []model.ScheduleParticipant `json:"participants"`
}

func (s *ScheduleService) Detail(id, tenantID uint) (*ScheduleDetail, error) {
	schedule, err := s.scheduleRepo.GetByID(id, tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("日程不存在")
	}
	participants, _ := s.scheduleRepo.ListParticipants(id, tenantID)
	return &ScheduleDetail{Schedule: *schedule, Participants: participants}, nil
}

func (s *ScheduleService) Create(tenantID, creatorID uint, title, description string, startTime, endTime time.Time, isAllDay, priority int8, location, color string, participantIDs []uint) (*model.Schedule, error) {
	schedule := model.Schedule{
		TenantID:    tenantID,
		Title:       title,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		IsAllDay:    isAllDay,
		Priority:    priority,
		Location:    location,
		Color:       color,
		CreatorID:   creatorID,
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txSchedRepo := repository.NewScheduleRepo(tx)
		if err := txSchedRepo.Create(&schedule); err != nil {
			return err
		}
		for _, pid := range participantIDs {
			p := &model.ScheduleParticipant{ScheduleID: schedule.ID, UserID: pid, TenantID: tenantID}
			if err := txSchedRepo.CreateParticipant(p); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, utils.ErrInternal("创建日程失败")
	}
	return &schedule, nil
}

func (s *ScheduleService) Update(id, tenantID uint, updates map[string]interface{}) error {
	affected, err := s.scheduleRepo.Update(id, tenantID, updates)
	if err != nil {
		return utils.ErrInternal("更新日程失败")
	}
	if affected == 0 {
		return utils.ErrNotFound("日程不存在")
	}
	return nil
}

func (s *ScheduleService) Delete(id, tenantID uint) error {
	if err := s.scheduleRepo.Delete(id, tenantID); err != nil {
		return utils.ErrInternal("删除失败")
	}
	return nil
}

func (s *ScheduleService) WeekList(userID, tenantID uint) ([]model.Schedule, error) {
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	schedules, err := s.scheduleRepo.ListByUserAndWeek(userID, tenantID, startOfWeek, endOfWeek)
	if err != nil {
		return nil, utils.ErrInternal("获取周日程失败")
	}
	return schedules, nil
}
