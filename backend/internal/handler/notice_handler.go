package handler

import (
	"net/http"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	noticeService *service.NoticeService
}

func NewNoticeHandler(noticeService *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{noticeService: noticeService}
}

func (h *NoticeHandler) List(c *gin.Context) {
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	noticeType := c.DefaultQuery("type", "")
	keyword := c.DefaultQuery("keyword", "")
	userID, _ := c.Get("user_id")

	var noticeTypeInt int
	if noticeType != "" {
		noticeTypeInt, _ = strconv.Atoi(noticeType)
	}
	list, total, err := h.noticeService.List(tid, noticeTypeInt, keyword, page, pageSize, userID.(uint))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": list, "total": total}})
}

func (h *NoticeHandler) Detail(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	notice, err := h.noticeService.Detail(uint(id), tid, userID.(uint))
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notice})
}

func (h *NoticeHandler) Create(c *gin.Context) {
	tid := getTenantID(c)
	userID, _ := c.Get("user_id")
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Type    *int8  `json:"type"`
		Summary string `json:"summary"`
		Cover   string `json:"cover"`
		IsTop   *int8  `json:"isTop"`
		Status  *int8  `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	noticeType := int8(1)
	if req.Type != nil {
		noticeType = *req.Type
	}
	isTop := int8(0)
	if req.IsTop != nil {
		isTop = *req.IsTop
	}
	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}
	notice, err := h.noticeService.Create(tid, userID.(uint), req.Title, req.Content, noticeType, req.Summary, req.Cover, isTop, status)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notice})
}

func (h *NoticeHandler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	count, err := h.noticeService.UnreadCount(userID.(uint), tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"count": count}})
}

func (h *NoticeHandler) MarkRead(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	if err := h.noticeService.MarkRead(uint(id), userID.(uint), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}
