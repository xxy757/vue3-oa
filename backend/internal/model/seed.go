package model

import (
	"oa-saas/internal/pkg/utils"
	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	var count int64
	db.Model(&Role{}).Count(&count)
	if count > 0 {
		return nil
	}

	adminRole := Role{
		Name:        "管理员",
		Code:        "admin",
		Description: "系统管理员",
		Permissions: StringArray{"*"},
		Status:      1,
	}
	db.Create(&adminRole)

	userRole := Role{
		Name:        "普通员工",
		Code:        "employee",
		Description: "普通员工角色",
		Permissions: StringArray{"approval:apply", "notice:view", "schedule:view"},
		Status:      1,
	}
	db.Create(&userRole)

	managerRole := Role{
		Name:        "部门经理",
		Code:        "manager",
		Description: "部门经理角色",
		Permissions: StringArray{"approval:apply", "approval:approve", "notice:view", "notice:create", "schedule:view", "schedule:manage"},
		Status:      1,
	}
	db.Create(&managerRole)

	techDept := Department{
		Name:   "技术部",
		Sort:   1,
		Leader: "管理员",
		Phone:  "13800000001",
		Email:  "tech@oa.com",
		Status: 1,
	}
	db.Create(&techDept)

	productDept := Department{
		Name:   "产品部",
		Sort:   2,
		Leader: "张三",
		Phone:  "13800000002",
		Email:  "product@oa.com",
		Status: 1,
	}
	db.Create(&productDept)

	hrDept := Department{
		Name:   "人事部",
		Sort:   3,
		Leader: "李四",
		Phone:  "13800000003",
		Email:  "hr@oa.com",
		Status: 1,
	}
	db.Create(&hrDept)

	adminPwd, _ := utils.HashPassword("123456")
	admin := User{
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
		{Username: "zhangsan", Password: userPwd, Nickname: "张三", Email: "zhangsan@oa.com", Phone: "13800000002", DeptID: &productDept.ID, RoleID: &managerRole.ID, Status: 1},
		{Username: "lisi", Password: userPwd, Nickname: "李四", Email: "lisi@oa.com", Phone: "13800000003", DeptID: &hrDept.ID, RoleID: &managerRole.ID, Status: 1},
		{Username: "wangwu", Password: userPwd, Nickname: "王五", Email: "wangwu@oa.com", Phone: "13800000004", DeptID: &techDept.ID, RoleID: &userRole.ID, Status: 1},
		{Username: "zhaoliu", Password: userPwd, Nickname: "赵六", Email: "zhaoliu@oa.com", Phone: "13800000005", DeptID: &techDept.ID, RoleID: &userRole.ID, Status: 1},
		{Username: "sunqi", Password: userPwd, Nickname: "孙七", Email: "sunqi@oa.com", Phone: "13800000006", DeptID: &productDept.ID, RoleID: &userRole.ID, Status: 1},
	}
	for _, u := range users {
		db.Create(&u)
	}

	flows := []ApprovalFlow{
		{
			Name: "请假审批流程",
			Code: "leave",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "人事审批", Type: "approval", Approver: []uint{3}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			Name: "报销审批流程",
			Code: "expense",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "财务审批", Type: "approval", Approver: []uint{1}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			Name: "加班审批流程",
			Code: "overtime",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "通知人事", Type: "notify", Approver: []uint{3}},
			},
			Status: 1,
		},
		{
			Name: "出差审批流程",
			Code: "travel",
			Nodes: FlowNodes{
				FlowNode{Name: "提交申请", Type: "submit", Approver: []uint{}},
				FlowNode{Name: "部门经理审批", Type: "approval", Approver: []uint{2}},
				FlowNode{Name: "人事审批", Type: "approval", Approver: []uint{3}},
				FlowNode{Name: "通知申请人", Type: "notify", Approver: []uint{}},
			},
			Status: 1,
		},
		{
			Name: "通用审批流程",
			Code: "general",
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
