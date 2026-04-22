package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepo
}

func NewRoleService(roleRepo *repository.RoleRepo) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) List(tenantID uint) ([]model.Role, error) {
	roles, err := s.roleRepo.ListByTenant(tenantID)
	if err != nil {
		return nil, utils.ErrInternal("获取角色列表失败")
	}
	return roles, nil
}

func (s *RoleService) Create(tenantID uint, name, code, description string, permissions []string, status int8) (*model.Role, error) {
	count, err := s.roleRepo.CountByCode(code, tenantID)
	if err != nil {
		return nil, utils.ErrInternal("查询角色失败")
	}
	if count > 0 {
		return nil, utils.ErrConflict("角色编码已存在")
	}
	role := model.Role{
		TenantID:    tenantID,
		Name:        name,
		Code:        code,
		Description: description,
		Permissions: permissions,
		Status:      status,
	}
	if err := s.roleRepo.Create(&role); err != nil {
		return nil, utils.ErrInternal("创建角色失败")
	}
	return &role, nil
}

func (s *RoleService) Update(id, tenantID uint, updates map[string]interface{}) error {
	affected, err := s.roleRepo.Update(id, tenantID, updates)
	if err != nil {
		return utils.ErrInternal("更新角色失败")
	}
	if affected == 0 {
		return utils.ErrNotFound("角色不存在")
	}
	return nil
}

func (s *RoleService) Delete(id, tenantID uint) error {
	role, err := s.roleRepo.GetByID(id, tenantID)
	if err != nil {
		return utils.ErrNotFound("角色不存在")
	}
	if role.Code == "admin" {
		return utils.ErrForbidden("无法删除管理员角色")
	}
	if err := s.roleRepo.Delete(role); err != nil {
		return utils.ErrInternal("删除失败")
	}
	return nil
}
