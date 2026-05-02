package models

import (
	"time"

	"gorm.io/gorm"
)

type Carousel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"not null"`
	ImageURL  string         `json:"image_url" gorm:"not null"`
	LinkURL   string         `json:"link_url"`
	Order     int            `json:"order" gorm:"default:0"`
	Status    string         `json:"status" gorm:"default:'active'"` // active, inactive
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
