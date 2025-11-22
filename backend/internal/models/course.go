package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title        string    `gorm:"size:255;not null"`
	Slug         string    `gorm:"size:255;uniqueIndex;not null"`
	Description  string    `gorm:"type:text"`
	PriceInPaisa int       `gorm:"default:0"`
	Topics       []Topic   `gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
