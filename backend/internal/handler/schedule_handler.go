package handler

import (
	"net/http"
	"oa-saas/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

func (h *ScheduleHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	startDate := c.DefaultQuery("startDate", "")
	endDate := c.DefaultQuery("endDate", "")
	schedules, err := h.scheduleService.List(userID.(uint), tid, startDate, endDate)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedules})
}

func (h *ScheduleHandler) Detail(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := h.scheduleService.Detail(uint(id), tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	var req struct {
		Title          string `json:"title" binding:"required"`
		Description    string `json:"description"`
		StartTime      string `json:"startTime" binding:"required"`
		EndTime        string `json:"endTime" binding:"required"`
		IsAllDay       *int8  `json:"isAllDay"`
		Priority       *int8  `json:"priority"`
		Location       string `json:"location"`
		Color          string `json:"color"`
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
	schedule, err := h.scheduleService.Create(tid, userID.(uint), req.Title, req.Description, startTime, endTime, isAllDay, priority, req.Location, color, req.ParticipantIDs)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedule})
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	tid := getTenantID(c)
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
	if err := h.scheduleService.Update(uint(id), tid, updates); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.scheduleService.Delete(uint(id), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}

func (h *ScheduleHandler) WeekList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	schedules, err := h.scheduleService.WeekList(userID.(uint), tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": schedules})
}
