package service

import (
	"oa-saas/internal/model"
	"oa-saas/internal/pkg/jwt"
	"oa-saas/internal/pkg/utils"
	"oa-saas/internal/repository"
)

type AuthService struct {
	userRepo   *repository.UserRepo
	deptRepo   *repository.DeptRepo
	roleRepo   *repository.RoleRepo
	tenantRepo *repository.TenantRepo
	planRepo   *repository.PlanRepo
	jwtSecret  string
	jwtExpire  int
}

func NewAuthService(
	userRepo *repository.UserRepo,
	deptRepo *repository.DeptRepo,
	roleRepo *repository.RoleRepo,
	tenantRepo *repository.TenantRepo,
	planRepo *repository.PlanRepo,
	jwtSecret string,
	jwtExpire int,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		deptRepo:   deptRepo,
		roleRepo:   roleRepo,
		tenantRepo: tenantRepo,
		planRepo:   planRepo,
		jwtSecret:  jwtSecret,
		jwtExpire:  jwtExpire,
	}
}

type LoginResult struct {
	Token string    `json:"token"`
	User  UserBasic `json:"user"`
	Tenant TenantBasic `json:"tenant"`
}

type UserBasic struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	DeptID   *uint  `json:"deptId"`
	RoleID   *uint  `json:"roleId"`
}

type TenantBasic struct {
	ID     uint     `json:"id"`
	Name   string   `json:"name"`
	Slug   string   `json:"slug"`
	Status string   `json:"status"`
	Plan   PlanBasic `json:"plan"`
}

type PlanBasic struct {
	ID       uint                `json:"id"`
	Name     string              `json:"name"`
	Code     string              `json:"code"`
	Features model.FeatureMap    `json:"features"`
	MaxUsers int                 `json:"maxUsers"`
}

type UserInfoResult struct {
	ID          uint              `json:"id"`
	Username    string            `json:"username"`
	Nickname    string            `json:"nickname"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	Avatar      string            `json:"avatar"`
	DeptID      *uint             `json:"deptId"`
	RoleID      *uint             `json:"roleId"`
	DeptName    string            `json:"deptName"`
	RoleName    string            `json:"roleName"`
	Permissions model.StringArray `json:"permissions"`
	CreateTime  interface{}       `json:"createTime"`
	Tenant      TenantBasic       `json:"tenant"`
}

func (s *AuthService) Login(username string, tenantID uint, password string, tenant model.Tenant) (*LoginResult, error) {
	user, err := s.userRepo.GetByUsername(username, tenantID)
	if err != nil {
		return nil, utils.ErrUnauthorized("用户名或密码错误")
	}
	if user.Status != 1 {
		return nil, utils.ErrForbidden("用户已被禁用")
	}
	if !utils.CheckPassword(password, user.Password) {
		return nil, utils.ErrUnauthorized("用户名或密码错误")
	}
	token, err := jwt.GenerateToken(user.ID, user.TenantID, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, utils.ErrInternal("生成token失败")
	}
	plan, _ := s.planRepo.GetByID(tenant.PlanID)
	result := &LoginResult{
		Token: token,
		User: UserBasic{
			ID:       user.ID,
			Username: user.Username,
			Nickname: user.Nickname,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			DeptID:   user.DeptID,
			RoleID:   user.RoleID,
		},
	}
	if plan != nil {
		result.Tenant = TenantBasic{
			ID:     tenant.ID,
			Name:   tenant.Name,
			Slug:   tenant.Slug,
			Status: tenant.Status,
			Plan: PlanBasic{
				ID:       plan.ID,
				Name:     plan.Name,
				Code:     plan.Code,
				Features: plan.Features,
				MaxUsers: plan.MaxUsers,
			},
		}
	}
	return result, nil
}

func (s *AuthService) ChangePassword(userID, tenantID uint, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(userID, tenantID)
	if err != nil {
		return utils.ErrNotFound("用户不存在")
	}
	if !utils.CheckPassword(oldPassword, user.Password) {
		return utils.ErrBadRequest("原密码错误")
	}
	hashedPwd, err := utils.HashPassword(newPassword)
	if err != nil {
		return utils.ErrInternal("密码加密失败")
	}
	return s.userRepo.UpdatePassword(userID, tenantID, hashedPwd)
}

func (s *AuthService) GetInfo(userID, tenantID uint, tenant model.Tenant) (*UserInfoResult, error) {
	user, err := s.userRepo.GetByID(userID, tenantID)
	if err != nil {
		return nil, utils.ErrNotFound("用户不存在")
	}
	result := &UserInfoResult{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Phone:     user.Phone,
		Avatar:    user.Avatar,
		DeptID:    user.DeptID,
		RoleID:    user.RoleID,
		CreateTime: user.CreatedAt,
	}
	if user.DeptID != nil {
		if dept, err := s.deptRepo.GetByID(*user.DeptID, tenantID); err == nil {
			result.DeptName = dept.Name
		}
	}
	if user.RoleID != nil {
		if role, err := s.roleRepo.GetByID(*user.RoleID, tenantID); err == nil {
			result.RoleName = role.Name
			result.Permissions = role.Permissions
		}
	}
	plan, _ := s.planRepo.GetByID(tenant.PlanID)
	if plan != nil {
		result.Tenant = TenantBasic{
			ID:     tenant.ID,
			Name:   tenant.Name,
			Slug:   tenant.Slug,
			Status: tenant.Status,
			Plan: PlanBasic{
				ID:       plan.ID,
				Name:     plan.Name,
				Code:     plan.Code,
				Features: plan.Features,
				MaxUsers: plan.MaxUsers,
			},
		}
	}
	return result, nil
}
