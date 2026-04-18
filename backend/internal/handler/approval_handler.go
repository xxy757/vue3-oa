package handler

import (
	"net/http"
	"oa-saas/internal/model"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ApprovalHandler 审批管理处理器，处理审批的创建、查询、操作和统计。
// 支持多节点审批流程（通过→下一节点→完成），以及驳回、转交、撤回等操作。
type ApprovalHandler struct {
	db *gorm.DB // 数据库实例
}

// NewApprovalHandler 创建审批管理处理器实例。
func NewApprovalHandler(db *gorm.DB) *ApprovalHandler {
	return &ApprovalHandler{db: db}
}

// Create 发起审批申请。
// 路由：POST /api/v1/approvals
//
// 执行步骤：
//  1. 获取当前用户ID和租户ID
//  2. 绑定并验证请求参数（title、type、content 必填）
//  3. 根据 type 查找对应的审批流程模板
//  4. 创建审批记录（状态为 pending）
//  5. 根据流程模板创建审批节点，第一个节点设为 active
//  6. 返回创建的审批数据
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

	// 查找对应的审批流程模板
	var flow model.ApprovalFlow
	if err := h.db.Where("code = ? AND tenant_id = ?", req.Type, tid).First(&flow).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到对应的审批流程"})
		return
	}

	// 创建审批记录
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

	// 根据流程模板创建审批节点
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
		// 第一个节点设为 active 状态
		if i == 0 {
			approvalNode.Status = "active"
		}
		h.db.Create(&approvalNode)
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": approval})
}

// MyList 获取当前用户发起的审批列表（我发起的）。
// 路由：GET /api/v1/approvals/my
//
// 执行步骤：
//  1. 获取当前用户ID和租户ID
//  2. 解析分页参数和筛选条件（type、status）
//  3. 按 applicant_id 过滤，支持可选的类型和状态筛选
//  4. 统计总数并执行分页查询
//  5. 返回分页列表数据
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

// PendingList 获取当前用户待审批的审批列表（待我审批）。
// 路由：GET /api/v1/approvals/pending
//
// 核心逻辑：通过 JOIN approval_nodes 表，使用 MySQL 的 JSON_CONTAINS 函数
// 匹配当前用户是否在审批人列表中，筛选节点状态为 active 或 pending 的记录。
func (h *ApprovalHandler) PendingList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var approvals []model.Approval
	var total int64
	query := h.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?",
			tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"active", "pending"})

	query.Count(&total)
	query.Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

// DoneList 获取当前用户已处理过的审批列表（我已审批）。
// 路由：GET /api/v1/approvals/done
//
// 核心逻辑：与 PendingList 类似，但筛选节点状态为 approved 或 rejected，
// 并使用 GROUP BY 去重（同一审批可能有多条匹配节点）。
func (h *ApprovalHandler) DoneList(c *gin.Context) {
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	var approvals []model.Approval
	var total int64
	query := h.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?",
			tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"approved", "rejected"})

	query.Count(&total)
	query.Group("approvals.id").Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"list": approvals, "total": total}})
}

// Detail 获取审批详情，包含审批基本信息和所有审批节点。
// 路由：GET /api/v1/approvals/:id
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

// Action 对审批执行操作（通过/驳回/转交）。
// 路由：POST /api/v1/approvals/:id/action
//
// 请求体：
//   - action：操作类型（approve/reject/transfer）
//   - comment：审批意见
//
// 执行步骤：
//  1. 查询审批记录，确认存在且状态为 pending
//  2. 查找当前活跃的审批节点
//  3. 根据操作类型执行不同逻辑：
//     a. approve：更新节点为 approved，查找下一节点并激活；若无下一节点则审批完成
//     b. reject：更新节点为 rejected，整个审批标记为 rejected
//     c. transfer：更新节点为 transferred，需额外提供转交目标用户ID
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

	// 查询审批记录
	var approval model.Approval
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).First(&approval).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审批不存在"})
		return
	}
	if approval.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "审批已结束"})
		return
	}

	// 查找当前活跃的审批节点
	var node model.ApprovalNode
	if err := h.db.Where("approval_id = ? AND tenant_id = ? AND status IN ?", id, tid, []string{"active", "pending"}).First(&node).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "未找到待审批节点"})
		return
	}

	now := time.Now()
	switch req.Action {
	case "approve":
		// 通过操作：更新当前节点状态
		node.Status = "approved"
		node.Comment = req.Comment
		node.ApproverID = userIDToUint(userID)
		node.UpdatedAt = now
		h.db.Save(&node)

		// 查找下一个审批节点
		var nextNode model.ApprovalNode
		if err := h.db.Where("approval_id = ? AND tenant_id = ? AND sort > ?", id, tid, node.Sort).Order("sort ASC").First(&nextNode).Error; err != nil {
			// 没有下一个节点，审批流程结束
			approval.Status = "approved"
			h.db.Save(&approval)
		} else {
			// 激活下一个节点
			nextNode.Status = "active"
			h.db.Save(&nextNode)
			approval.CurrentNode = nextNode.Sort
			h.db.Save(&approval)
		}
	case "reject":
		// 驳回操作：标记节点和整个审批为已驳回
		node.Status = "rejected"
		node.Comment = req.Comment
		node.ApproverID = userIDToUint(userID)
		node.UpdatedAt = now
		h.db.Save(&node)
		approval.Status = "rejected"
		h.db.Save(&approval)
	case "transfer":
		// 转交操作：需要额外提供目标用户ID
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

// Withdraw 撤回审批申请。仅申请人本人可以撤回状态为 pending 的审批。
// 路由：POST /api/v1/approvals/:id/withdraw
func (h *ApprovalHandler) Withdraw(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	userID, _ := c.Get("user_id")
	tid := getTenantID(c)

	var approval model.Approval
	if err := h.db.Where("id = ? AND tenant_id = ?", id, tid).First(&approval).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "审批不存在"})
		return
	}

	// 仅申请人可撤回
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

// Stats 获取当前用户的审批统计数据。
// 路由：GET /api/v1/approvals/stats
//
// 返回数据：
//   - myPending：我发起的待审批数
//   - myApproved：我发起的已通过数
//   - myRejected：我发起的已驳回数
//   - pendingApproval：待我审批的数量
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

	// 统计待我审批的数量
	var pendingApproval int64
	h.db.Model(&model.ApprovalNode{}).
		Where("tenant_id = ? AND JSON_CONTAINS(approver_ids, ?) AND status IN ?",
			tid, strconv.FormatUint(uint64(userID.(uint)), 10), []string{"active", "pending"}).
		Count(&pendingApproval)

	c.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{
		"myPending":       pending,
		"myApproved":      approved,
		"myRejected":      rejected,
		"pendingApproval": pendingApproval,
	}})
}

// userIDToUint 将接口类型的用户ID转换为 uint 类型。
// 用于从 gin.Context 的 Get 方法返回值中安全地提取用户ID。
func userIDToUint(id interface{}) uint {
	if v, ok := id.(uint); ok {
		return v
	}
	return 0
}
