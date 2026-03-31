package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FlowHandler struct {
	db *gorm.DB
}

func NewFlowHandler(db *gorm.DB) *FlowHandler {
	return &FlowHandler{db: db}
}

func (h *FlowHandler) List(c *gin.Context) {
	var flows []model.ApprovalFlow
	if err := h.db.Find(&flows).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "获取流程列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": flows})
}

func (h *FlowHandler) Create(c *gin.Context) {
	var req struct {
		Name        string           `json:"name" binding:"required"`
		Code        string           `json:"code" binding:"required"`
		Description string           `json:"description"`
		Nodes       []model.FlowNode `json:"nodes" binding:"required"`
		Status      *int8            `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var count int64
	h.db.Model(&model.ApprovalFlow{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"code": 409, "message": "流程编码已存在"})
		return
	}

	status := int8(1)
	if req.Status != nil {
		status = *req.Status
	}

	flow := model.ApprovalFlow{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Nodes:       req.Nodes,
		Status:      status,
	}
	if err := h.db.Create(&flow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建流程失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": flow})
}

func (h *FlowHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		Name        string           `json:"name" binding:"required"`
		Description string           `json:"description"`
		Nodes       []model.FlowNode `json:"nodes" binding:"required"`
		Status      *int8            `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"nodes":       req.Nodes,
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	result := h.db.Model(&model.ApprovalFlow{}).Where("id = ?", id).Updates(updates)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "流程不存在"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "更新成功"})
}

func (h *FlowHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.db.Delete(&model.ApprovalFlow{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "删除成功"})
}
