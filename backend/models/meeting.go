package models

import (
	"time"

	"gorm.io/gorm"
)

type Meeting struct {
	ID            uint              `json:"id" gorm:"primaryKey"`
	Title         string            `json:"title" gorm:"not null"`
	Description   string            `json:"description"`
	MeetingRoomID uint              `json:"meeting_room_id" gorm:"not null"`
	Room          MeetingRoom       `json:"room" gorm:"foreignKey:MeetingRoomID"`
	OrganizerID   uint              `json:"organizer_id" gorm:"not null"`
	Organizer     User              `json:"organizer" gorm:"foreignKey:OrganizerID"`
	StartTime     time.Time         `json:"start_time" gorm:"not null"`
	EndTime       time.Time         `json:"end_time" gorm:"not null"`
	Status        string            `json:"status" gorm:"default:'scheduled'"` // scheduled, ongoing, completed, cancelled
	Reminders     []MeetingReminder `json:"reminders" gorm:"foreignKey:MeetingID"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	DeletedAt     gorm.DeletedAt    `json:"-" gorm:"index"`
}
