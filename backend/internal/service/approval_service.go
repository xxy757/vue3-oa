package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ApprovalService struct {
	approvalRepo *repository.ApprovalRepo
	flowRepo     *repository.FlowRepo
	db           *gorm.DB
}

func NewApprovalService(approvalRepo *repository.ApprovalRepo, flowRepo *repository.FlowRepo, db *gorm.DB) *ApprovalService {
	return &ApprovalService{
		approvalRepo: approvalRepo,
		flowRepo:     flowRepo,
		db:           db,
	}
}

func (s *ApprovalService) Create(tenantID, applicantID uint, title, approvalType string, content model.JSONObject) (*model.Approval, error) {
	flow, err := s.flowRepo.GetByCode(approvalType, tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("未找到对应的审批流程")
	}
	approval := &model.Approval{
		TenantID:    tenantID,
		Title:       title,
		Type:        approvalType,
		Content:     content,
		ApplicantID: applicantID,
		Status:      "pending",
		CurrentNode: 0,
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txApprovalRepo := repository.NewApprovalRepo(tx)
		if err := txApprovalRepo.Create(approval); err != nil {
			return err
		}
		for i, node := range flow.Nodes {
			approverIDs := make(model.UintArray, 0)
			for _, id := range node.Approver {
				approverIDs = append(approverIDs, id)
			}
			approvalNode := model.ApprovalNode{
				TenantID:    tenantID,
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
			if err := txApprovalRepo.CreateNode(&approvalNode); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, utils.ErrInternal("创建审批失败")
	}
	return approval, nil
}

func (s *ApprovalService) MyList(userID, tenantID uint, approvalType, status string, page, pageSize int) ([]model.Approval, int64, error) {
	return s.approvalRepo.MyList(userID, tenantID, approvalType, status, page, pageSize)
}

func (s *ApprovalService) PendingList(userID, tenantID uint, page, pageSize int) ([]model.Approval, int64, error) {
	return s.approvalRepo.PendingList(userID, tenantID, page, pageSize)
}

func (s *ApprovalService) DoneList(userID, tenantID uint, page, pageSize int) ([]model.Approval, int64, error) {
	return s.approvalRepo.DoneList(userID, tenantID, page, pageSize)
}

type ApprovalDetail struct {
	Approval model.Approval     `json:"approval"`
	Nodes    []model.ApprovalNode `json:"nodes"`
}

func (s *ApprovalService) Detail(id, tenantID uint) (*ApprovalDetail, error) {
	approval, err := s.approvalRepo.GetByID(id, tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("审批不存在")
	}
	nodes, _ := s.approvalRepo.ListNodes(id, tenantID)
	return &ApprovalDetail{Approval: *approval, Nodes: nodes}, nil
}

func (s *ApprovalService) Action(id, userID, tenantID uint, action, comment string) error {
	approval, err := s.approvalRepo.GetByID(id, tenantID)
	if err != nil {
		return utils.ErrNotFound("审批不存在")
	}
	if approval.Status != "pending" {
		return utils.ErrBadRequest("审批已结束")
	}
	node, err := s.approvalRepo.GetActiveNode(id, tenantID)
	if err != nil {
		return utils.ErrNotFound("未找到待审批节点")
	}

	now := time.Now()
	switch action {
	case "approve":
		node.Status = "approved"
		node.Comment = comment
		node.ApproverID = userID
		node.UpdatedAt = now
		s.approvalRepo.SaveNode(node)
		nextNode, err := s.approvalRepo.GetNextNode(id, tenantID, node.Sort)
		if err != nil {
			approval.Status = "approved"
			s.approvalRepo.Save(approval)
		} else {
			nextNode.Status = "active"
			s.approvalRepo.SaveNode(nextNode)
			approval.CurrentNode = nextNode.Sort
			s.approvalRepo.Save(approval)
		}
	case "reject":
		node.Status = "rejected"
		node.Comment = comment
		node.ApproverID = userID
		node.UpdatedAt = now
		s.approvalRepo.SaveNode(node)
		approval.Status = "rejected"
		s.approvalRepo.Save(approval)
	case "transfer":
		node.Status = "transferred"
		node.Comment = comment
		node.ApproverID = userID
		s.approvalRepo.SaveNode(node)
	default:
		return utils.ErrBadRequest("不支持的操作")
	}
	return nil
}

func (s *ApprovalService) Withdraw(id, userID, tenantID uint) error {
	approval, err := s.approvalRepo.GetByID(id, tenantID)
	if err != nil {
		return utils.ErrNotFound("审批不存在")
	}
	if approval.ApplicantID != userID {
		return utils.ErrForbidden("无权操作")
	}
	if approval.Status != "pending" {
		return utils.ErrBadRequest("只能撤回待审批的申请")
	}
	approval.Status = "withdrawn"
	return s.approvalRepo.Save(approval)
}

type ApprovalStats struct {
	MyPending       int64 `json:"myPending"`
	MyApproved      int64 `json:"myApproved"`
	MyRejected      int64 `json:"myRejected"`
	PendingApproval int64 `json:"pendingApproval"`
}

func (s *ApprovalService) Stats(userID, tenantID uint) (*ApprovalStats, error) {
	stats, err := s.approvalRepo.Stats(userID, tenantID)
	if err != nil {
		return nil, utils.ErrInternal("获取统计失败")
	}
	result := &ApprovalStats{}
	for _, st := range stats {
		switch st.Status {
		case "pending":
			result.MyPending = st.Count
		case "approved":
			result.MyApproved = st.Count
		case "rejected":
			result.MyRejected = st.Count
		}
	}
	pendingApproval, _ := s.approvalRepo.CountPendingApprovalNodes(userID, tenantID)
	result.PendingApproval = pendingApproval
	return result, nil
}
