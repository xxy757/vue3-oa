package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
)

type FlowService struct {
	flowRepo *repository.FlowRepo
}

func NewFlowService(flowRepo *repository.FlowRepo) *FlowService {
	return &FlowService{flowRepo: flowRepo}
}

func (s *FlowService) List(tenantID uint) ([]model.ApprovalFlow, error) {
	flows, err := s.flowRepo.ListByTenant(tenantID)
	if err != nil {
		return nil, utils.ErrInternal("获取流程列表失败")
	}
	return flows, nil
}

func (s *FlowService) Create(tenantID uint, name, code, description string, nodes []model.FlowNode, status int8) (*model.ApprovalFlow, error) {
	count, err := s.flowRepo.CountByCode(code, tenantID)
	if err != nil {
		return nil, utils.ErrInternal("查询流程失败")
	}
	if count > 0 {
		return nil, utils.ErrConflict("流程编码已存在")
	}
	flow := model.ApprovalFlow{
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Description: description,
		Nodes:       nodes,
		Status:      status,
	}
	if err := s.flowRepo.Create(&flow); err != nil {
		return nil, utils.ErrInternal("创建流程失败")
	}
	return &flow, nil
}

func (s *FlowService) Update(id, tenantID uint, updates map[string]interface{}) error {
	affected, err := s.flowRepo.Update(id, tenantID, updates)
	if err != nil {
		return utils.ErrInternal("更新流程失败")
	}
	if affected == 0 {
		return utils.ErrNotFound("流程不存在")
	}
	return nil
}

func (s *FlowService) Delete(id, tenantID uint) error {
	if err := s.flowRepo.Delete(id, tenantID); err != nil {
		return utils.ErrInternal("删除失败")
	}
	return nil
}
