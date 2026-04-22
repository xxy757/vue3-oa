package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"

	"gorm.io/gorm"
)

type UserService struct {
	userRepo   *repository.UserRepo
	tenantRepo *repository.TenantRepo
	deptRepo   *repository.DeptRepo
	roleRepo   *repository.RoleRepo
	db         *gorm.DB
}

func NewUserService(userRepo *repository.UserRepo, tenantRepo *repository.TenantRepo, deptRepo *repository.DeptRepo, roleRepo *repository.RoleRepo, db *gorm.DB) *UserService {
	return &UserService{
		userRepo:   userRepo,
		tenantRepo: tenantRepo,
		deptRepo:   deptRepo,
		roleRepo:   roleRepo,
		db:         db,
	}
}

type UserItem struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Avatar     string `json:"avatar"`
	DeptID     *uint  `json:"deptId"`
	RoleID     *uint  `json:"roleId"`
	DeptName   string `json:"deptName"`
	RoleName   string `json:"roleName"`
	Status     int8   `json:"status"`
	CreateTime string `json:"createTime"`
}

func (s *UserService) List(tenantID uint, keyword string, page, pageSize int) ([]UserItem, int64, error) {
	users, total, err := s.userRepo.ListByTenant(tenantID, keyword, page, pageSize)
	if err != nil {
		return nil, 0, utils.ErrInternal("获取用户列表失败")
	}
	var list []UserItem
	for _, u := range users {
		item := UserItem{
			ID:         u.ID,
			Username:   u.Username,
			Nickname:   u.Nickname,
			Email:      u.Email,
			Phone:      u.Phone,
			Avatar:     u.Avatar,
			DeptID:     u.DeptID,
			RoleID:     u.RoleID,
			Status:     u.Status,
			CreateTime: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if u.DeptID != nil {
			if dept, err := s.deptRepo.GetByID(*u.DeptID, tenantID); err == nil {
				item.DeptName = dept.Name
			}
		}
		if u.RoleID != nil {
			if role, err := s.roleRepo.GetByID(*u.RoleID, tenantID); err == nil {
				item.RoleName = role.Name
			}
		}
		list = append(list, item)
	}
	return list, total, nil
}

func (s *UserService) Create(tenantID uint, username, password, nickname, email, phone string, deptID, roleID *uint, status int8) (*model.User, error) {
	count, err := s.userRepo.CountByUsername(username, tenantID)
	if err != nil {
		return nil, utils.ErrInternal("查询用户失败")
	}
	if count > 0 {
		return nil, utils.ErrConflict("用户名已存在")
	}
	tenant, err := s.tenantRepo.GetByID(tenantID)
	if err == nil && tenant.CurrentUsers >= tenant.MaxUsers {
		return nil, utils.ErrForbidden("套餐用户数已满，请升级套餐")
	}
	hashedPwd, err := utils.HashPassword(password)
	if err != nil {
		return nil, utils.ErrInternal("密码加密失败")
	}
	user := model.User{
		TenantID: tenantID,
		Username: username,
		Password: hashedPwd,
		Nickname: nickname,
		Email:    email,
		Phone:    phone,
		DeptID:   deptID,
		RoleID:   roleID,
		Status:   status,
	}
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		txUserRepo := repository.NewUserRepo(tx)
		txTenantRepo := repository.NewTenantRepo(tx)
		if err := txUserRepo.Create(&user); err != nil {
			return err
		}
		return txTenantRepo.IncrementCurrentUsers(tenantID)
	}); err != nil {
		return nil, utils.ErrInternal("创建用户失败")
	}
	return &user, nil
}

func (s *UserService) Update(id, tenantID uint, updates map[string]interface{}) error {
	affected, err := s.userRepo.Update(id, tenantID, updates)
	if err != nil {
		return utils.ErrInternal("更新用户失败")
	}
	if affected == 0 {
		return utils.ErrNotFound("用户不存在")
	}
	return nil
}

func (s *UserService) Delete(id, tenantID uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		txUserRepo := repository.NewUserRepo(tx)
		txTenantRepo := repository.NewTenantRepo(tx)
		if err := txUserRepo.Delete(id, tenantID); err != nil {
			return err
		}
		return txTenantRepo.DecrementCurrentUsers(tenantID)
	})
}

func (s *UserService) UpdateStatus(id, tenantID uint, status int8) error {
	return s.userRepo.UpdateStatus(id, tenantID, status)
}
