package model

import "time"

type Department struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"index;not null" json:"tenantId"`
	ParentID  *uint     `gorm:"index" json:"parentId"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Leader    string    `gorm:"size:50" json:"leader"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Email     string    `gorm:"size:100" json:"email"`
	Status    int8      `gorm:"default:1;not null" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
