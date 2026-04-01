package model

import "time"

type Notice struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	TenantID    uint      `gorm:"index;not null" json:"tenantId"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	Type        int8      `gorm:"default:1" json:"type"`
	Summary     string    `gorm:"size:500" json:"summary"`
	Cover       string    `gorm:"size:255" json:"cover"`
	IsTop       int8      `gorm:"default:0" json:"isTop"`
	Status      int8      `gorm:"default:1" json:"status"`
	PublisherID uint      `gorm:"index;not null" json:"publisherId"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type NoticeRead struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TenantID  uint      `gorm:"index;not null" json:"tenantId"`
	NoticeID  uint      `gorm:"uniqueIndex:idx_notice_user;not null" json:"noticeId"`
	UserID    uint      `gorm:"uniqueIndex:idx_notice_user;not null" json:"userId"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createTime"`
}
