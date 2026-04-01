package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON: %v", value)
	}
	return json.Unmarshal(bytes, s)
}

type Role struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	TenantID    uint        `gorm:"index;not null" json:"tenantId"`
	Name        string      `gorm:"size:50;not null" json:"name"`
	Code        string      `gorm:"uniqueIndex:uk_tenant_code;size:50;not null" json:"code"`
	Description string      `gorm:"size:255" json:"description"`
	Permissions StringArray `gorm:"type:json" json:"permissions"`
	Status      int8        `gorm:"default:1;not null" json:"status"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updateTime"`
}
