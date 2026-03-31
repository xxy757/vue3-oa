package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type FlowNodes []FlowNode

func (f FlowNodes) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FlowNodes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON: %v", value)
	}
	return json.Unmarshal(bytes, f)
}

type FlowNode struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	Approver []uint `json:"approver"`
}

type ApprovalFlow struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	Nodes       FlowNodes `gorm:"type:json;not null" json:"nodes"`
	Status      int8      `gorm:"default:1;not null" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}
