package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
)

type DeptService struct {
	deptRepo *repository.DeptRepo
}

func NewDeptService(deptRepo *repository.DeptRepo) *DeptService {
	return &DeptService{deptRepo: deptRepo}
}

type DeptTreeNode struct {
	ID       uint            `json:"id"`
	ParentID *uint           `json:"parentId"`
	Name     string          `json:"name"`
	Sort     int             `json:"sort"`
	Leader   string          `json:"leader"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	Status   int8            `json:"status"`
	Children []*DeptTreeNode `json:"children,omitempty"`
}

func (s *DeptService) List(tenantID uint) ([]*DeptTreeNode, error) {
	depts, err := s.deptRepo.ListByTenant(tenantID)
	if err != nil {
		return nil, utils.ErrInternal("获取部门列表失败")
	}
	return buildDeptTree(depts, nil), nil
}

func buildDeptTree(depts []model.Department, parentID *uint) []*DeptTreeNode {
	var nodes []*DeptTreeNode
	for i := range depts {
		if (depts[i].ParentID == nil && parentID == nil) || (depts[i].ParentID != nil && parentID != nil && *depts[i].ParentID == *parentID) {
			node := &DeptTreeNode{
				ID:       depts[i].ID,
				ParentID: depts[i].ParentID,
				Name:     depts[i].Name,
				Sort:     depts[i].Sort,
				Leader:   depts[i].Leader,
				Phone:    depts[i].Phone,
				Email:    depts[i].Email,
				Status:   depts[i].Status,
			}
			node.Children = buildDeptTree(depts, &depts[i].ID)
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *DeptService) Create(tenantID uint, name string, parentID *uint, sort int, leader, phone, email string, status int8) (*model.Department, error) {
	dept := model.Department{
		TenantID: tenantID,
		ParentID: parentID,
		Name:     name,
		Sort:     sort,
		Leader:   leader,
		Phone:    phone,
		Email:    email,
		Status:   status,
	}
	if err := s.deptRepo.Create(&dept); err != nil {
		return nil, utils.ErrInternal("创建部门失败")
	}
	return &dept, nil
}

func (s *DeptService) Update(id, tenantID uint, updates map[string]interface{}) error {
	affected, err := s.deptRepo.Update(id, tenantID, updates)
	if err != nil {
		return utils.ErrInternal("更新部门失败")
	}
	if affected == 0 {
		return utils.ErrNotFound("部门不存在")
	}
	return nil
}

func (s *DeptService) Delete(id, tenantID uint) error {
	count, err := s.deptRepo.CountChildren(id, tenantID)
	if err != nil {
		return utils.ErrInternal("查询子部门失败")
	}
	if count > 0 {
		return utils.ErrConflict("该部门下有子部门，无法删除")
	}
	if err := s.deptRepo.Delete(id, tenantID); err != nil {
		return utils.ErrInternal("删除失败")
	}
	return nil
}
