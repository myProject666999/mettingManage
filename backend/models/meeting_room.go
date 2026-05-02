package models

import (
	"time"

	"gorm.io/gorm"
)

type MeetingRoom struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Location    string         `json:"location"`
	Capacity    int            `json:"capacity" gorm:"not null"`
	Description string         `json:"description"`
	Equipment   string         `json:"equipment"` // 设备，用逗号分隔
	Status      string         `json:"status" gorm:"default:'available'"` // available, unavailable
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
