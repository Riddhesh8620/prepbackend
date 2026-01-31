package handlers

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePaymentSession(c *fiber.Ctx) error {

	tx := store.DB.Begin()

	user_id, _ := uuid.Parse(c.Locals("user_id").(string))
	payment_id, _ := uuid.Parse(c.FormValue("payment_id"))
	payment_mode := c.FormValue("payment_mode")
	payable_amount, _ := strconv.ParseFloat(c.FormValue("amount"), 32)

	// 1. Build Model
	sessionModel := models.UserPaymentSession{
		UserID:        user_id,
		PaymentID:     payment_id,
		PaymentMode:   payment_mode,
		PayableAmount: float32(payable_amount),
	}

	if err := tx.Create(&sessionModel).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Transaction failed"})
	}

	return c.Status(201).JSON(&sessionModel)
}
