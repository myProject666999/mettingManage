package models

import (
	"time"

	"gorm.io/gorm"
)

type MeetingDocument struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	MeetingID   uint           `json:"meeting_id" gorm:"not null"`
	Meeting     Meeting        `json:"meeting" gorm:"foreignKey:MeetingID"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	FileType    string         `json:"file_type"` // pdf, doc, docx, ppt, etc.
	FileSize    int64          `json:"file_size"`
	FilePath    string         `json:"file_path" gorm:"not null"`
	UploaderID  uint           `json:"uploader_id" gorm:"not null"`
	Uploader    interface{}    `json:"uploader" gorm:"-"` // 可以是User或Admin
	DownloadCount int          `json:"download_count" gorm:"default:0"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
