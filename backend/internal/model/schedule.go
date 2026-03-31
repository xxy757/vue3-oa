package model

import "time"

type Schedule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	StartTime   time.Time `gorm:"not null;index:idx_time" json:"startTime"`
	EndTime     time.Time `gorm:"not null;index:idx_time" json:"endTime"`
	IsAllDay    int8      `gorm:"default:0" json:"isAllDay"`
	Priority    int8      `gorm:"default:1" json:"priority"`
	Location    string    `gorm:"size:255" json:"location"`
	Color       string    `gorm:"size:20;default:'#1677FF'" json:"color"`
	CreatorID   uint      `gorm:"index;not null" json:"creatorId"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"createTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updateTime"`
}

type ScheduleParticipant struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ScheduleID uint      `gorm:"uniqueIndex:idx_schedule_user;not null" json:"scheduleId"`
	UserID     uint      `gorm:"uniqueIndex:idx_schedule_user;not null" json:"userId"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createTime"`
}
