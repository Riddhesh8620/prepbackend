package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base defines the common columns for all tables
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // json:"-" hides it from API responses
}

// This function runs automatically before any record is inserted into the DB.
func (base *Base) BeforeCreate(tx *gorm.DB) (err error) {
	uuidv7, err := uuid.NewV7()
	if err == nil && base.ID == uuid.Nil {
		base.ID = uuidv7
	}
	return
}
