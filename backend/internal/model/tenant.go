package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

type FeatureMap map[string]interface{}

func (f FeatureMap) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FeatureMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSON: %v", value)
	}
	return json.Unmarshal(bytes, f)
}

type Plan struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:50;not null" json:"name"`
	Code      string     `gorm:"uniqueIndex;size:30;not null" json:"code"`
	Price     float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	MinUsers  int        `gorm:"default:1" json:"minUsers"`
	MaxUsers  int        `gorm:"default:100" json:"maxUsers"`
	Features  FeatureMap `gorm:"type:json" json:"features"`
	IsActive  int8       `gorm:"default:1" json:"isActive"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
}

type Tenant struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Slug         string     `gorm:"uniqueIndex;size:50;not null" json:"slug"`
	Logo         string     `gorm:"size:255" json:"logo"`
	ContactName  string     `gorm:"size:50;not null" json:"contactName"`
	ContactPhone string     `gorm:"size:20;not null" json:"contactPhone"`
	ContactEmail string     `gorm:"size:100;not null" json:"contactEmail"`
	PlanID       uint       `gorm:"index;not null" json:"planId"`
	PlanStartAt  *time.Time `json:"planStartAt"`
	PlanExpireAt *time.Time `json:"planExpireAt"`
	CurrentUsers int        `gorm:"default:0" json:"currentUsers"`
	MaxUsers     int        `gorm:"default:5" json:"maxUsers"`
	Status       string     `gorm:"size:20;default:trial;index" json:"status"`
	TrialEndsAt  *time.Time `json:"trialEndsAt"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
	DeletedAt    *time.Time `gorm:"index" json:"-"`
}

type Invoice struct {
	ID                   uint       `gorm:"primaryKey" json:"id"`
	TenantID             uint       `gorm:"index;not null" json:"tenantId"`
	PlanID               uint       `gorm:"not null" json:"planId"`
	InvoiceNo            string     `gorm:"uniqueIndex;size:50;not null" json:"invoiceNo"`
	PeriodStart          time.Time  `gorm:"not null" json:"periodStart"`
	PeriodEnd            time.Time  `gorm:"not null" json:"periodEnd"`
	UserCount            int        `gorm:"not null" json:"userCount"`
	Amount               float64    `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status               string     `gorm:"size:20;default:pending;index" json:"status"`
	PaidAt               *time.Time `json:"paidAt"`
	PaymentMethod        string     `gorm:"size:20" json:"paymentMethod"`
	PaymentTransactionID string     `gorm:"size:100" json:"paymentTransactionId"`
	CreatedAt            time.Time  `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt            time.Time  `gorm:"autoUpdateTime" json:"updateTime"`
}

type TenantLog struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	TenantID   uint       `gorm:"index;not null" json:"tenantId"`
	OperatorID *uint      `json:"operatorId"`
	Action     string     `gorm:"size:50;index;not null" json:"action"`
	TargetType string     `gorm:"size:50" json:"targetType"`
	TargetID   *uint      `json:"targetId"`
	Content    JSONObject `gorm:"type:json" json:"content"`
	IP         string     `gorm:"size:45" json:"ip"`
	UserAgent  string     `gorm:"size:500" json:"userAgent"`
	CreatedAt  time.Time  `gorm:"autoCreateTime;index" json:"createTime"`
}
