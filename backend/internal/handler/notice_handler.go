package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NoticeHandler struct {
	db *gorm.DB
}

func NewNoticeHandler(db *gorm.DB) *NoticeHandler {
	return &NoticeHandler{db: db}
}

func (h *NoticeHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	noticeType := c.DefaultQuery("type", "")
	keyword := c.DefaultQuery("keyword", "")

	var notices []model.Notice
	var total int64
	query := h.db.Model(&model.Notice{})
	if noticeType != "" {
		t, _ := strconv.Atoi(noticeType)
		query = query.Where("type = ?", t)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}
	query.Count(&total)
	query.Order("is_top DESC, created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&notices)

	userID, _ := c.Get("user_id")
	type NoticeItem struct {
		model.Notice
		PublisherName string `json:"publisherName"`
		IsRead        bool   `json:"isRead"`
	}
	var list []NoticeItem
	for _, n := range notices {
		item := NoticeItem{Notice: n}
		var user model.User
		if err := h.db.First(&user, n.PublisherID).Error; err == nil {
			item.PublisherName = user.Nickname
		}
		var read model.NoticeRead
		result := h.db.Where("notice_id = ? AND user_id = ?", n.ID, userID).First(&read)
		item.IsRead = result.RowsAffected > 0
		list = append(list, item)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": list, "total": total}})
}

func (h *NoticeHandler) Detail(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var notice model.Notice
	if err := h.db.First(&notice, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "公告不存在"})
		return
	}
	userID, _ := c.Get("user_id")
	var read model.NoticeRead
	h.db.Where("notice_id = ? AND user_id = ?", id, userID).FirstOrCreate(&read, model.NoticeRead{NoticeID: uint(id), UserID: userID.(uint)})
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notice})
}

func (h *NoticeHandler) Create(c *gin.Context) {
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

	notice := model.Notice{
		Title:       req.Title,
		Content:     req.Content,
		Type:        noticeType,
		Summary:     req.Summary,
		Cover:       req.Cover,
		IsTop:       isTop,
		Status:      status,
		PublisherID: userID.(uint),
	}
	if err := h.db.Create(&notice).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "发布公告失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": notice})
}

func (h *NoticeHandler) UnreadCount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var count int64
	h.db.Model(&model.Notice{}).
		Where("id NOT IN (SELECT notice_id FROM notice_reads WHERE user_id = ?)", userID).
		Count(&count)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"count": count}})
}

func (h *NoticeHandler) MarkRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	var read model.NoticeRead
	h.db.Where("notice_id = ? AND user_id = ?", id, userID).FirstOrCreate(&read, model.NoticeRead{NoticeID: uint(id), UserID: userID.(uint)})
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "标记成功"})
}
