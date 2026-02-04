package models

import (
	"time"

	"github.com/google/uuid"
)

type TopicInventory struct {
	BaseID
	UserID               uuid.UUID `gorm:"type:uuid" json:"user_id"`
	IsActive             bool      `json:"is_active"`
	TopicID              uuid.UUID `gorm:"type:uuid" json:"topic_id"`
	UserPaymentSessionId uuid.UUID `gorm:"type:uuid" json:"user_payment_session_id"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`

	Topic   Topic              `gorm:"foreignKey:TopicID"`
	Session UserPaymentSession `gorm:"foreignKey:UserPaymentSessionId"`
	User    User               `gorm:"foreignKey:UserID"`
}

type CourseInventory struct {
	BaseID
	UserID               uuid.UUID `gorm:"type:uuid" json:"user_id"`
	IsActive             bool      `json:"is_active"`
	CourseID             uuid.UUID `gorm:"type:uuid" json:"course_id"`
	UserPaymentSessionId uuid.UUID `gorm:"type:uuid" json:"user_payment_session_id"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`

	Course  Course             `gorm:"foreignKey:CourseID"`
	Session UserPaymentSession `gorm:"foreignKey:UserPaymentSessionId"`
	User    User               `gorm:"foreignKey:UserID"`
}

type AddInventoryReqDto struct {
	ID               string    `json:"product_id"`
	Type             string    `json:"type"` // Topic | Course
	Price            float32   `json:"product_amount"`
	PaymentSessionID string    `json:"payment_session_id"`
	CreationDateTime time.Time `json:"creation_datetime"`
}

type TopicGroup struct {
	TopicTitle  string `gorm:"column:topic_title"`
	CourseTitle string `gorm:"column:course_title"`
}

type CourseGroup struct {
	Topics []string
}

type CourseProject struct {
	CourseTitle string `gorm:"column:course_title"`
}
