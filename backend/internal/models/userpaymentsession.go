package models

import (
	"time"

	"github.com/google/uuid"
)

type UserPaymentSession struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID        uuid.UUID `gorm:"type:uuid" json:"user_id"`
	PaymentID     string    `gorm:"type:varchar;"`
	PaymentMode   string    `gorm:"size:100" json:"payment_mode"`
	PayableAmount float32   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Status        string    `gorm:"size:100" json:"status"`
	CreatedOn     time.Time `json:"created_on"`
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

type PaymentSessionCreateDto struct {
	PaymentSessionID string               `json:"payment_session_id"`
	CreationDateTime time.Time            `json:"creation_datetime"`
	CartTotal        float32              `json:"cart_total"`
	TransactionId    uuid.UUID            `json:"transaction_id"`
	Payload          []AddInventoryReqDto `json:"payload"`
}

type PaymentSessionResponse struct {
	ID               string  `json:"id"`
	UserID           string  `json:"user_id"`
	PaymentSessionID string  `json:"payment_session_id"`
	PaymentMode      string  `json:"payment_mode"`
	CartTotal        float32 `json:"cart_total"`
	Status           string  `json:"status"`
}

type UpdatePaymentRequest struct {
	Status string `json:"status"` // "Completed" or "Failed"
}

type GetPaymentDashboardDto struct {
	ID        string  `json:"id"`
	UserEmail string  `json:"user_email"`
	Name      string  `json:"name"`
	Amount    float32 `json:"amount"`
	Status    string  `json:"status"`
}
