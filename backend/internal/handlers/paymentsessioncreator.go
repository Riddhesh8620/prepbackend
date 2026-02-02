package handlers

import (
	"errors"
	"net/http"
	"prepbackend/internal/models"
	"prepbackend/internal/store"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreatePaymentSession(c *fiber.Ctx) error {
	var sessionRequest models.PaymentSessionCreateDto

	tx := store.DB.Begin()

	if err := c.BodyParser(&sessionRequest); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	var inventoryTopic []models.TopicInventory
	var inventoryCourse []models.CourseInventory

	user_id, _ := uuid.Parse(c.Locals("user_id").(string))
	// 1. Build Model
	sessionModel := models.UserPaymentSession{
		ID:            uuid.New(),
		UserID:        user_id,
		PaymentID:     sessionRequest.PaymentSessionID,
		PaymentMode:   "WhatsApp",
		PayableAmount: float32(sessionRequest.CartTotal),
		Status:        models.PaymentStatus.Verifying,
	}
	// 1. Creating a session for payment, which will be in verifying stage acknowledged by System Administrator.
	if err := tx.Create(&sessionModel).Error; err != nil {
		tx.Rollback()
		return c.Status(400).JSON(fiber.Map{"error": "Transaction failed"})
	}

	// 2. Creating inventory for current user in in-active state, activated on successful payment.
	for _, item := range sessionRequest.Payload {
		product_id, _ := uuid.Parse(item.ID)

		switch item.Type {
		case "course":

			inventoryCourse = append(inventoryCourse, models.CourseInventory{
				UserID:               user_id,
				CourseID:             product_id,
				UserPaymentSessionId: sessionModel.ID, // Link to the payment session
				IsActive:             false,           // In-active until Admin approves
				CreatedAt:            time.Now().UTC(),
				UpdatedAt:            time.Now().UTC(),
			})

		case "topic":
			inventoryTopic = append(inventoryTopic, models.TopicInventory{
				UserID:               user_id,
				TopicID:              product_id,
				UserPaymentSessionId: sessionModel.ID,
				IsActive:             false,
				CreatedAt:            time.Now().UTC(),
				UpdatedAt:            time.Now().UTC(),
			})
		}
	}

	tx.Create(&inventoryCourse)

	tx.Create(&inventoryTopic)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.Status(400).JSON(fiber.Map{"error": "Failed to initiate transaction, pls do not panic if money has been debited from your account, contact System Admin!"})
	}

	response := models.PaymentSessionResponse{
		ID:               sessionModel.ID.String(),
		UserID:           user_id.String(),
		PaymentSessionID: sessionModel.PaymentID,
		PaymentMode:      sessionModel.PaymentMode,
		CartTotal:        sessionModel.PayableAmount,
		Status:           models.PaymentStatus.Verifying,
	}

	return c.Status(201).JSON(&response)
}

func HandlePaymentSessionHook(c *fiber.Ctx) error {
	session_id, err := uuid.Parse(c.Params("session_id"))

	if err != nil || session_id == uuid.Nil {
		return c.SendStatus(http.StatusNotFound)
	}

	var session models.UserPaymentSession
	result := store.DB.Where("id = ?", session_id).First(&session)

	if result.Error != nil {
		// Handle record not found specifically
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Payment session not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error occurred",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":     session.ID,
		"status": session.Status, // Assuming your model has a 'Status' field
	})
}

func HandleAdminUpdatePayment(c *fiber.Ctx) error {
	// 1. Parse Session ID from URL
	sessionID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid session ID format",
		})
	}

	// 2. Parse Request Body
	var req models.UpdatePaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// 3. Start Database Transaction
	err = store.DB.Transaction(func(tx *gorm.DB) error {
		// A. Update the User Payment Session Status
		// We use a map to ensure GORM doesn't skip fields
		if err := tx.Model(&models.UserPaymentSession{}).
			Where("id = ?", sessionID).
			Update("status", req.Status).Error; err != nil {
			return err
		}

		// B. If Status is Success, Unlock Inventories
		if req.Status == "SUCCESS" { // Or models.PaymentStatus.Success

			// Unlock Course Inventory
			if err := tx.Model(&models.CourseInventory{}).
				Where("user_payment_session_id = ?", sessionID).
				Update("is_active", true).Error; err != nil {
				return err
			}

			// Unlock Topic Inventory
			if err := tx.Model(&models.TopicInventory{}).
				Where("user_payment_session_id = ?", sessionID).
				Update("is_active", true).Error; err != nil {
				return err
			}
		}

		// Return nil to commit the transaction
		return nil
	})

	// 4. Handle Transaction Outcome
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update payment records: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment status updated and inventory unlocked successfully",
		"status":  req.Status,
	})
}
