package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"index;not null" json:"tenantId"`
	Username  string    `gorm:"uniqueIndex:uk_tenant_username;size:50;not null" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Nickname  string    `gorm:"size:50;not null" json:"nickname"`
	Email     string    `gorm:"size:100" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Avatar    string    `gorm:"size:255" json:"avatar"`
	DeptID    *uint     `gorm:"index" json:"deptId"`
	RoleID    *uint     `gorm:"index" json:"roleId"`
	Status    int8      `gorm:"default:1;not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
