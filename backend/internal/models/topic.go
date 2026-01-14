package models

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	Base
	CourseID    uuid.UUID `gorm:"type:uuid;index"`
	Title       string    `gorm:"size:255;not null"`
	Description string    `gorm:"type:text"`
	OrderIndex  int
	Price       float32   `gorm:"type:decimal(10,2);default:10"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsDeleted   bool      `gorm:"default:false"`
	DeletedAt   time.Time `gorm:"default:NULL"`
}
