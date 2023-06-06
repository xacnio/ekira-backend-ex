package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID        uuid.UUID      `gorm:"column:session_id;primary_key;type:varchar(36)" json:"session_id"`
	UserID    uuid.UUID      `gorm:"column:user_id;type:varchar(36);not null" json:"user_id"`
	UserAgent string         `gorm:"column:user_agent;type:varchar(255);not null" json:"user_agent"`
	IP        string         `gorm:"column:ip;type:varchar(255);not null" json:"ip"`
	ExpiresAt int64          `gorm:"column:expires_at;not null" json:"expires_at"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at" validate:""`
}
