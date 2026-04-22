package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"oa-saas/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApprovalHandler struct {
	approvalService *service.ApprovalService
}

func NewApprovalHandler(approvalService *service.ApprovalService) *ApprovalHandler {
	return &ApprovalHandler{approvalService: approvalService}
}

func (h *ApprovalHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	var req struct {
		Title   string           `json:"title" binding:"required"`
		Type    string           `json:"type" binding:"required"`
		Content model.JSONObject `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	approval, err := h.approvalService.Create(tid, userID.(uint), req.Title, req.Type, req.Content)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": approval})
}

func (h *ApprovalHandler) MyList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	approvalType := c.DefaultQuery("type", "")
	status := c.DefaultQuery("status", "")
	approvals, total, err := h.approvalService.MyList(userID.(uint), tid, approvalType, status, page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) PendingList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	approvals, total, err := h.approvalService.PendingList(userID.(uint), tid, page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) DoneList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	approvals, total, err := h.approvalService.DoneList(userID.(uint), tid, page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) Detail(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	result, err := h.approvalService.Detail(uint(id), tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}

func (h *ApprovalHandler) Action(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	var req struct {
		Action  string `json:"action" binding:"required"`
		Comment string `json:"comment"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}
	if err := h.approvalService.Action(uint(id), userID.(uint), tid, req.Action, req.Comment); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

func (h *ApprovalHandler) Withdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	if err := h.approvalService.Withdraw(uint(id), userID.(uint), tid); err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "撤回成功"})
}

func (h *ApprovalHandler) Stats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	result, err := h.approvalService.Stats(userID.(uint), tid)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": result})
}
