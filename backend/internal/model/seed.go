package model

import (
	"oa-saas/internal/pkg/utils"
	"time"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	var count int64
	db.Model(&Role{}).Count(&count)
	if count > 0 {
		return nil
	}

	plans := []Plan{
		{Name: "免费版", Code: "free", Price: 0, MinUsers: 1, MaxUsers: 5, Features: FeatureMap{"approval": false, "schedule": false, "notice": true, "storage": 100}, IsActive: 1},
		{Name: "标准版", Code: "standard", Price: 29, MinUsers: 5, MaxUsers: 50, Features: FeatureMap{"approval": true, "schedule": true, "notice": true, "storage": 1024}, IsActive: 1},
		{Name: "专业版", Code: "professional", Price: 59, MinUsers: 10, MaxUsers: 200, Features: FeatureMap{"approval": true, "schedule": true, "notice": true, "storage": 5120, "api": true}, IsActive: 1},
		{Name: "企业版", Code: "enterprise", Price: 99, MinUsers: 50, MaxUsers: 1000, Features: FeatureMap{"approval": true, "schedule": true, "notice": true, "storage": 20480, "api": true, "sso": true}, IsActive: 1},
	}
	for i := range plans {
		db.Create(&plans[i])
	}

	trialEnds := time.Now().Add(14 * 24 * time.Hour)
	tenant := Tenant{
		Name:         "示例科技有限公司",
		Slug:         "demo",
		ContactName:  "管理员",
		ContactPhone: "13800000001",
		ContactEmail: "admin@demo.com",
		PlanID:       plans[1].ID,
		MaxUsers:     plans[1].MaxUsers,
		Status:       "trial",
		TrialEndsAt:  &trialEnds,
	}
	db.Create(&tenant)

	adminRole := Role{
		TenantID:    tenant.ID,
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员",
		Permissions: StringArray{"*"},
		Status:      1,
	}
	db.Create(&adminRole)

	userRole := Role{
		TenantID:    tenant.ID,
		Name:        "普通员工",
		Code:        "employee",
		Description: "普通员工角色",
		Permissions: StringArray{"approval:apply", "notice:view", "schedule:view"},
		Status:      1,
	}
	db.Create(&userRole)

	managerRole := Role{
		TenantID:    tenant.ID,
		Name:        "部门经理",
		Code:        "manager",
		Description: "部门经理角色",
		Permissions: StringArray{"approval:apply", "approval:approve", "notice:view", "notice:create", "schedule:view", "schedule:manage"},
		Status:      1,
	}
	db.Create(&managerRole)

	techDept := Department{
		TenantID: tenant.ID,
		Name:     "技术部",
		Sort:     1,
		Leader:   "管理员",
		Phone:    "13800000001",
		Email:    "tech@oa.com",
		Status:   1,
	}
	db.Create(&techDept)

	productDept := Department{
		TenantID: tenant.ID,
		Name:     "产品部",
		Sort:     2,
		Leader:   "张三",
		Phone:    "13800000002",
		Email:    "product@oa.com",
		Status:   1,
	}
	db.Create(&productDept)

	hrDept := Department{
		TenantID: tenant.ID,
		Name:     "人事部",
		Sort:     3,
		Leader:   "李四",
		Phone:    "13800000003",
		Email:    "hr@oa.com",
		Status:   1,
	}
	db.Create(&hrDept)

	adminPwd, _ := utils.HashPassword("123456")
	admin := User{
		TenantID: tenant.ID,
		Username: "admin",
		Password: adminPwd,
		Nickname: "管理员",
		Email:    "admin@oa.com",
		Phone:    "13800000001",
		DeptID:   &techDept.ID,
		RoleID:   &adminRole.ID,
		Status:   1,
	}
	db.Create(&admin)

	userPwd, _ := utils.HashPassword("123456")
	users := []User{
		{TenantID: tenant.ID, Username: "zhangsan", Password: userPwd, Nickname: "张三", Email: "zhangsan@oa.com", Phone: "13800000002", DeptID: &productDept.ID, RoleID: &managerRole.ID, Status: 1},
		{TenantID: tenant.ID, Username: "lisi", Password: userPwd, Nickname: "李四", Email: "lisi@oa.com", Phone: "13800000003", DeptID: &hrDept.ID, RoleID: &managerRole.ID, Status: 1},
		{TenantID: tenant.ID, Username: "wangwu", Password: userPwd, Nickname: "王五", Email: "wangwu@oa.com", Phone: "13800000004", DeptID: &techDept.ID, RoleID: &userRole.ID, Status: 1},
		{TenantID: tenant.ID, Username: "zhaoliu", Password: userPwd, Nickname: "赵六", Email: "zhaoliu@oa.com", Phone: "13800000005", DeptID: &techDept.ID, RoleID: &userRole.ID, Status: 1},
		{TenantID: tenant.ID, Username: "sunqi", Password: userPwd, Nickname: "孙七", Email: "sunqi@oa.com", Phone: "13800000006", DeptID: &productDept.ID, RoleID: &userRole.ID, Status: 1},
	}
	for _, u := range users {
		db.Create(&u)
	}

	db.Model(&Tenant{}).Where("id = ?", tenant.ID).Update("current_users", 6)

	flows := []ApprovalFlow{
		{
			TenantID: tenant.ID,
			Name:     "请假审批流程",
			Code:     "leave",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "人事审批", Type: "approval", Approver: []uint{3}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			TenantID: tenant.ID,
			Name:     "报销审批流程",
			Code:     "expense",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "财务审批", Type: "approval", Approver: []uint{1}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			TenantID: tenant.ID,
			Name:     "加班审批流程",
			Code:     "overtime",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "通知人事", Type: "notify", Approver: []uint{3}},
			},
			Status: 1,
		},
		{
			TenantID: tenant.ID,
			Name:     "出差审批流程",
			Code:     "travel",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "人事审批", Type: "approval", Approver: []uint{3}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			TenantID: tenant.ID,
			Name:     "通用审批流程",
			Code:     "general",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
			},
			Status: 1,
		},
	}
	for _, f := range flows {
		db.Create(&f)
	}

	return nil
}
