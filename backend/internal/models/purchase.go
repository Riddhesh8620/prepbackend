package models

import (
	"github.com/google/uuid"
)

type Purchase struct {
	Base
	UserID    uuid.UUID  `gorm:"type:uuid;index"`
	TopicID   *uuid.UUID `gorm:"type:uuid;index;default:NULL"`
	CourseID  *uuid.UUID `gorm:"type:uuid;index;default:NULL"`
	Amount    float32    `gorm:"type:decimal(10,2);not null"`
	PaymentID string     `gorm:"size:255"`
	Status    string     `gorm:"size:50;default:'pending'"`
}
