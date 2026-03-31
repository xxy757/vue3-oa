package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	db *gorm.DB
}

func NewScheduleHandler(db *gorm.DB) *ScheduleHandler {
	return &ScheduleHandler{db: db}
}

func (h *ScheduleHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")

	var schedules []model.Schedule
	query := h.db.Model(&model.Schedule{}).
		Joins("LEFT JOIN schedule_participants ON schedule_participants.schedule_id = schedules.id").
		Where("schedules.creator_id = ? OR schedule_participants.user_id = ?", userID, userID)
	if startDate != "" && endDate != "" {
		start, _ := time.Parse("2006-01-02", startDate)
		end, _ := time.Parse("2006-01-02", endDate)
		query = query.Where("schedules.start_time <= ? AND schedules.end_time >= ?", end, start)
	}
	query.Distinct("schedules.id").Order("schedules.start_time ASC").Find(&schedules)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedules})
}

func (h *ScheduleHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var schedule model.Schedule
	if err := h.db.First(&schedule, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "日程不存在"})
		return
	}
	var participants []model.ScheduleParticipant
	h.db.Where("schedule_id = ?", id).Find(&participants)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"schedule": schedule, "participants": participants}})
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req struct {
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description"`
		StartTime     string `json:"startTime" binding:"required"`
		EndTime       string `json:"endTime" binding:"required"`
		IsAllDay      *int8  `json:"isAllDay"`
		Priority      *int8  `json:"priority"`
		Location      string `json:"location"`
		Color         string `json:"color"`
		ParticipantIDs []uint `json:"participantIds"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	isAllDay := int8(0)
	if req.IsAllDay != nil {
		isAllDay = *req.IsAllDay
	}
	priority := int8(1)
	if req.Priority != nil {
		priority = *req.Priority
	}
	color := "#1677FF"
	if req.Color != "" {
		color = req.Color
	}

	schedule := model.Schedule{
		Title:       req.Title,
		Description: req.Description,
		StartTime:   startTime,
		EndTime:     endTime,
		IsAllDay:    isAllDay,
		Priority:    priority,
		Location:    req.Location,
		Color:       color,
		CreatorID:   userID.(uint),
	}
	if err := h.db.Create(&schedule).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建日程失败"})
		return
	}

	for _, pid := range req.ParticipantIDs {
		h.db.Create(&model.ScheduleParticipant{ScheduleID: schedule.ID, UserID: pid})
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedule})
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartTime   string `json:"startTime"`
		EndTime     string `json:"endTime"`
		IsAllDay    *int8  `json:"isAllDay"`
		Priority    *int8  `json:"priority"`
		Location    string `json:"location"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartTime != "" {
		t, _ := time.Parse(time.RFC3339, req.StartTime)
		updates["start_time"] = t
	}
	if req.EndTime != "" {
		t, _ := time.Parse(time.RFC3339, req.EndTime)
		updates["end_time"] = t
	}
	if req.IsAllDay != nil {
		updates["is_all_day"] = *req.IsAllDay
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Location != "" {
		updates["location"] = req.Location
	}
	if req.Color != "" {
		updates["color"] = req.Color
	}

	result := h.db.Model(&model.Schedule{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "日程不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.db.Where("schedule_id = ?", id).Delete(&model.ScheduleParticipant{})
	if err := h.db.Delete(&model.Schedule{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ScheduleHandler) WeekList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	now := time.Now()
	startOfWeek := now.AddDate(0, 0, -int(now.Weekday()))
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	var schedules []model.Schedule
	h.db.Model(&model.Schedule{}).
		Joins("LEFT JOIN schedule_participants ON schedule_participants.schedule_id = schedules.id").
		Where("(schedules.creator_id = ? OR schedule_participants.user_id = ?) AND schedules.start_time >= ? AND schedules.start_time < ?", userID, userID, startOfWeek, endOfWeek).
		Distinct("schedules.id").
		Order("schedules.start_time ASC").
		Find(&schedules)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedules})
}
