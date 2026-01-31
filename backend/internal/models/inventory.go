package models

import (
	"time"

	"github.com/google/uuid"
)

type TopicInventory struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID               uuid.UUID `json:"user_id"`
	IsActive             bool      `json:"is_active"`
	TopicID              uuid.UUID `json:"topic_id"`
	UserPaymentSessionId uuid.UUID `json:"user_payment_session_id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CourseInventory struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	CourseID  uuid.UUID `json:"course_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddInventoryReqDto struct {
	Type string `json:"type"`
}
