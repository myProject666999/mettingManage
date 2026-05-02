package models

import (
	"time"

	"gorm.io/gorm"
)

type MeetingReminder struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MeetingID   uint           `json:"meeting_id" gorm:"not null"`
	Meeting     Meeting        `json:"meeting" gorm:"foreignKey:MeetingID"`
	ReminderType string        `json:"reminder_type" gorm:"not null"` // email, sms, app
	ReminderTime time.Time     `json:"reminder_time" gorm:"not null"`
	Message     string         `json:"message"`
	IsSent      bool           `json:"is_sent" gorm:"default:false"`
	SentTime    *time.Time     `json:"sent_time"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
