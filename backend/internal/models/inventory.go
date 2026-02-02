package models

import (
	"time"

	"github.com/google/uuid"
)

type TopicInventory struct {
	BaseID
	UserID               uuid.UUID `json:"user_id"`
	IsActive             bool      `json:"is_active"`
	TopicID              uuid.UUID `json:"topic_id"`
	UserPaymentSessionId uuid.UUID `json:"user_payment_session_id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

type CourseInventory struct {
	BaseID
	UserID               uuid.UUID `json:"user_id"`
	IsActive             bool      `json:"is_active"`
	CourseID             uuid.UUID `json:"course_id"`
	UserPaymentSessionId uuid.UUID `json:"user_payment_session_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type AddInventoryReqDto struct {
	ID               string    `json:"product_id"`
	Type             string    `json:"type"` // Topic | Course
	Price            float32   `json:"product_amount"`
	PaymentSessionID string    `json:"payment_session_id"`
	CreationDateTime time.Time `json:"creation_datetime"`
}
