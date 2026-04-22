package repository

import (
	"oa-saas/internal/model"
	"strconv"

	"gorm.io/gorm"
)

type ApprovalRepo struct {
	db *gorm.DB
}

func NewApprovalRepo(db *gorm.DB) *ApprovalRepo {
	return &ApprovalRepo{db: db}
}

func (r *ApprovalRepo) Create(approval *model.Approval) error {
	return r.db.Create(approval).Error
}

func (r *ApprovalRepo) CreateNode(node *model.ApprovalNode) error {
	return r.db.Create(node).Error
}

func (r *ApprovalRepo) GetByID(id uint, tenantID uint) (*model.Approval, error) {
	var approval model.Approval
	if err := r.db.Where("id = ? AND tenant_id = ?", id, tenantID).First(&approval).Error; err != nil {
		return nil, err
	}
	return &approval, nil
}

func (r *ApprovalRepo) Save(approval *model.Approval) error {
	return r.db.Save(approval).Error
}

func (r *ApprovalRepo) SaveNode(node *model.ApprovalNode) error {
	return r.db.Save(node).Error
}

func (r *ApprovalRepo) ListNodes(approvalID, tenantID uint) ([]model.ApprovalNode, error) {
	var nodes []model.ApprovalNode
	if err := r.db.Where("approval_id = ? AND tenant_id = ?", approvalID, tenantID).Order("sort ASC").Find(&nodes).Error; err != nil {
		return nil, err
	}
	return nodes, nil
}

func (r *ApprovalRepo) GetActiveNode(approvalID, tenantID uint) (*model.ApprovalNode, error) {
	var node model.ApprovalNode
	if err := r.db.Where("approval_id = ? AND tenant_id = ? AND status IN ?", approvalID, tenantID, []string{"active", "pending"}).First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *ApprovalRepo) GetNextNode(approvalID, tenantID uint, afterSort int) (*model.ApprovalNode, error) {
	var node model.ApprovalNode
	if err := r.db.Where("approval_id = ? AND tenant_id = ? AND sort > ?", approvalID, tenantID, afterSort).Order("sort ASC").First(&node).Error; err != nil {
		return nil, err
	}
	return &node, nil
}

func (r *ApprovalRepo) MyList(userID, tenantID uint, approvalType, status string, page, pageSize int) ([]model.Approval, int64, error) {
	var approvals []model.Approval
	var total int64
	query := r.db.Model(&model.Approval{}).Where("applicant_id = ? AND tenant_id = ?", userID, tenantID)
	if approvalType != "" {
		query = query.Where("type = ?", approvalType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals).Error; err != nil {
		return nil, 0, err
	}
	return approvals, total, nil
}

func (r *ApprovalRepo) PendingList(userID, tenantID uint, page, pageSize int) ([]model.Approval, int64, error) {
	var approvals []model.Approval
	var total int64
	query := r.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?",
			tenantID, strconv.FormatUint(uint64(userID), 10), []string{"active", "pending"})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals).Error; err != nil {
		return nil, 0, err
	}
	return approvals, total, nil
}

func (r *ApprovalRepo) DoneList(userID, tenantID uint, page, pageSize int) ([]model.Approval, int64, error) {
	var approvals []model.Approval
	var total int64
	query := r.db.Model(&model.Approval{}).
		Joins("JOIN approval_nodes ON approval_nodes.approval_id = approvals.id").
		Where("approvals.tenant_id = ? AND JSON_CONTAINS(approval_nodes.approver_ids, ?) AND approval_nodes.status IN ?",
			tenantID, strconv.FormatUint(uint64(userID), 10), []string{"approved", "rejected"})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Group("approvals.id").Order("approvals.created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&approvals).Error; err != nil {
		return nil, 0, err
	}
	return approvals, total, nil
}

func (r *ApprovalRepo) Stats(userID, tenantID uint) ([]struct {
	Status string
	Count  int64
}, error) {
	var stats []struct {
		Status string
		Count  int64
	}
	err := r.db.Model(&model.Approval{}).Select("status, count(*) as count").Where("applicant_id = ? AND tenant_id = ?", userID, tenantID).Group("status").Find(&stats).Error
	return stats, err
}

func (r *ApprovalRepo) CountPendingApprovalNodes(userID, tenantID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.ApprovalNode{}).
		Where("tenant_id = ? AND JSON_CONTAINS(approver_ids, ?) AND status IN ?",
			tenantID, strconv.FormatUint(uint64(userID), 10), []string{"active", "pending"}).
		Count(&count).Error
	return count, err
}
