package models

import (
	"time"

	"github.com/google/uuid"
)

type Topic struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CourseID     uuid.UUID `gorm:"type:uuid;index"`
	Title        string    `gorm:"size:255;not null"`
	Description  string    `gorm:"type:text"`
	OrderIndex   int
	PriceInPaisa int `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
