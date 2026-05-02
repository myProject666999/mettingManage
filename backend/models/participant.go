package models

import (
	"time"

	"gorm.io/gorm"
)

type Participant struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	MeetingID uint           `json:"meeting_id" gorm:"not null"`
	Meeting   Meeting        `json:"meeting" gorm:"foreignKey:MeetingID"`
	UserID    uint           `json:"user_id" gorm:"not null"`
	User      User           `json:"user" gorm:"foreignKey:UserID"`
	Status    string         `json:"status" gorm:"default:'invited'"` // invited, confirmed, declined, attended
	JoinTime  *time.Time     `json:"join_time"`
	LeaveTime *time.Time     `json:"leave_time"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
