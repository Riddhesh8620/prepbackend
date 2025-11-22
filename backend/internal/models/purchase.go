package models

import (
	"time"

	"github.com/google/uuid"
)

type Purchase struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID        uuid.UUID  `gorm:"type:uuid;index"`
	TopicID       *uuid.UUID `gorm:"type:uuid;index;default:NULL"`
	CourseID      *uuid.UUID `gorm:"type:uuid;index;default:NULL"`
	AmountInPaisa int
	PaymentID     string `gorm:"size:255"`
	Status        string `gorm:"size:50;default:'pending'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
