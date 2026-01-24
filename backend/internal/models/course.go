package models

import (
	"github.com/google/uuid"
)

type Course struct {
	Title         string  `gorm:"size:150;not null" json:"title"`
	Description   string  `gorm:"type:text" json:"description"`
	IsActive      bool    `gorm:"default:true" json:"is_active"`
	Thumbnail     string  `gorm:"type:text" json:"image"`
	Price         float32 `gorm:"not null" json:"price"`
	OriginalPrice float32 `gorm:"not null" json:"originalPrice"`
	Level         string  `gorm:"size:50;not null" json:"level"`
	Duration      string  `gorm:"not null" json:"duration"` // in minutes
	// Relations
	CategoryID uuid.UUID `json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Topics     []Topic   `json:"topics"`
	Base
}
