package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type JSONObject map[string]interface{}

func (j JSONObject) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONObject) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON: %v", value)
	}
	return json.Unmarshal(bytes, j)
}

type UintArray []uint

func (u UintArray) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UintArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON: %v", value)
	}
	return json.Unmarshal(bytes, u)
}

type Approval struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	TenantID    uint       `gorm:"index;not null" json:"tenantId"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Type        string     `gorm:"size:20;not null" json:"type"`
	Content     JSONObject `gorm:"type:json" json:"content"`
	ApplicantID uint       `gorm:"index;not null" json:"applicantId"`
	Status      string     `gorm:"size:20;default:pending;index" json:"status"`
	CurrentNode int        `gorm:"default:0" json:"currentNode"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
}

type ApprovalNode struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TenantID    uint      `gorm:"index;not null" json:"tenantId"`
	ApprovalID  uint      `gorm:"index;not null" json:"approvalId"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Type        string    `gorm:"size:20;not null" json:"type"`
	ApproverIDs UintArray `gorm:"type:json" json:"approverIds"`
	Status      string    `gorm:"size:20;default:pending" json:"status"`
	Comment     string    `gorm:"size:500" json:"comment"`
	ApproverID  uint      `gorm:"index" json:"approverId"`
	Sort        int       `gorm:"default:0" json:"sort"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
