package models

import (
	"github.com/google/uuid"
)

type UserPaymentSession struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid" json:"user_id"`
	PaymentID     uuid.UUID `gorm:"type:uuid;"`
	PaymentMode   string    `gorm:"size:100"`
	PayableAmount float32   `gorm:"type:decimal(10,2);not null"`
	Status        string    `gorm:"size:100"`
}

var PaymentStatus = struct {
	Pending   string
	Verifying string
	Success   string
	Failed    string
}{
	Pending:   "PENDING",
	Verifying: "VERIFYING",
	Success:   "SUCCESS",
	Failed:    "FAILED",
}
