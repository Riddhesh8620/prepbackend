package models

import (
	"github.com/google/uuid"
)

type Topic struct {
	CourseID    uuid.UUID `gorm:"type:uuid;index" json:"courseId"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	OrderIndex  int8      `gorm:"not null;default:0" json:"order"`
	Price       float32   `gorm:"type:decimal(10,2);default:10" json:"price"`
	Duration    string    `gorm:"size:100" json:"duration"` // in minutes
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	Base
}
