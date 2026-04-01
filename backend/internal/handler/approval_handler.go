package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApprovalHandler struct {
	db *gorm.DB
}

func NewApprovalHandler(db *gorm.DB) *ApprovalHandler {
	return &ApprovalHandler{db: db}
}

func (h *ApprovalHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	var req struct {
		Title   string                 `json:"title" binding:"required"`
		Type    string                 `json:"type" binding:"required"`
		Content map[string]interface{} `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	var flow model.ApprovalFlow
	if err := h.db.Where("code = ? AND tenant_id = ?", req.Type, tid).First(&flow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到对应的审批流程"})
		return
	}

	approval := model.Approval{
		TenantID:    tid,
		Title:       req.Title,
		Type:        req.Type,
		Content:     req.Content,
		ApplicantID: userID.(uint),
		Status:      "pending",
		CurrentNode: 0,
	}
	if err := h.db.Create(&approval).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "创建审批失败"})
		return
	}

	for i, node := range flow.Nodes {
		approverIDs := make(model.UintArray, 0)
		for _, id := range node.Approver {
			approverIDs = append(approverIDs, id)
		}
		approvalNode := model.ApprovalNode{
			TenantID:    tid,
			ApprovalID:  approval.ID,
			Name:        node.Name,
			Type:        node.Type,
			ApproverIDs: approverIDs,
			Status:      "pending",
			Sort:        i,
		}
		if i == 0 {
			approvalNode.Status = "active"
		}
		h.db.Create(&approvalNode)
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

	var approvals []model.Approval
	var total int64
	query := h.db.Model(&model.Approval{}).Where("applicant_id = ? AND tenant_id = ?", userID, tid)
	if approvalType != "" {
		query = query.Where("type = ?", approvalType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) PendingList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var approvals []model.Approval
	var total int64
	query := h.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?", tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"active", "pending"})
	query.Count(&total)
	query.Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) DoneList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var approvals []model.Approval
	var total int64
	query := h.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?", tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"approved", "rejected"})
	query.Count(&total)
	query.Group("approvals.id").Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

func (h *ApprovalHandler) Detail(c *gin.Context) {
	tid := getTenantID(c)
	id, _ := strconv.Atoi(c.Param("id"))
	var approval model.Approval
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).First(&approval).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审批不存在"})
		return
	}
	var nodes []model.ApprovalNode
	h.db.Where("approval_id = ? AND tenant_id = ?", id, tid).Order("sort ASC").Find(&nodes)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"approval": approval, "nodes": nodes}})
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

	var approval model.Approval
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).First(&approval).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审批不存在"})
		return
	}
	if approval.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "审批已结束"})
		return
	}

	var node model.ApprovalNode
	if err := h.db.Where("approval_id = ? AND tenant_id = ? AND status IN ?", id, tid, []string{"active", "pending"}).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到待审批节点"})
		return
	}

	now := time.Now()
	switch req.Action {
	case "approve":
		node.Status = "approved"
		node.Comment = req.Comment
		node.ApproverID = userIDToUint(userID)
		node.UpdatedAt = now
		h.db.Save(&node)
		var nextNode model.ApprovalNode
		if err := h.db.Where("approval_id = ? AND tenant_id = ? AND sort > ?", id, tid, node.Sort).Order("sort ASC").First(&nextNode).Error; err != nil {
			approval.Status = "approved"
			h.db.Save(&approval)
		} else {
			nextNode.Status = "active"
			h.db.Save(&nextNode)
			approval.CurrentNode = nextNode.Sort
			h.db.Save(&approval)
		}
	case "reject":
		node.Status = "rejected"
		node.Comment = req.Comment
		node.ApproverID = userIDToUint(userID)
		node.UpdatedAt = now
		h.db.Save(&node)
		approval.Status = "rejected"
		h.db.Save(&approval)
	case "transfer":
		var transferReq struct {
			TargetUserID uint `json:"targetUserId" binding:"required"`
		}
		if err := c.ShouldBindJSON(&transferReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "请指定转交人"})
			return
		}
		node.Status = "transferred"
		node.Comment = req.Comment
		node.ApproverID = userIDToUint(userID)
		h.db.Save(&node)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "不支持的操作"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "操作成功"})
}

func (h *ApprovalHandler) Withdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	var approval model.Approval
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).First(&approval).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审批不存在"})
		return
	}
	if approval.ApplicantID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "无权操作"})
		return
	}
	if approval.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "只能撤回待审批的申请"})
		return
	}
	approval.Status = "withdrawn"
	h.db.Save(&approval)
	c.JSON(http.StatusOK, gin.H{"code": 200, "message": "撤回成功"})
}

func (h *ApprovalHandler) Stats(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)
	type StatItem struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var stats []StatItem
	h.db.Model(&model.Approval{}).Select("status, count(*) as count").Where("applicant_id = ? AND tenant_id = ?", userID, tid).Group("status").Find(&stats)

	pending, approved, rejected := int64(0), int64(0), int64(0)
	for _, s := range stats {
		switch s.Status {
		case "pending":
			pending = s.Count
		case "approved":
			approved = s.Count
		case "rejected":
			rejected = s.Count
		}
	}

	var pendingApproval int64
	h.db.Model(&model.ApprovalNode{}).
		Where("tenant_id = ? AND JSON_CONTAINS(approver_ids, ?) AND status IN ?", tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"active", "pending"}).
		Count(&pendingApproval)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
		"myPending":       pending,
		"myApproved":      approved,
		"myRejected":      rejected,
		"pendingApproval": pendingApproval,
	}})
}

func userIDToUint(id interface{}) uint {
	if v, ok := id.(uint); ok {
		return v
	}
	return 0
}
