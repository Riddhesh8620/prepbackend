package models

import (
	"github.com/google/uuid"
)

type Course struct {
	Base
	Title       string `gorm:"size:150;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Relations
	CategoryID uuid.UUID `json:"category_id"`
	Category   Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Topics     []Topic   `json:"topics"`
}
