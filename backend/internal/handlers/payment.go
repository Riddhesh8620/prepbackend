package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"

	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type verifyReq struct {
	RazorpayPaymentID string `json:"razorpay_payment_id"`
	RazorpayOrderID   string `json:"razorpay_order_id"`
	RazorpaySignature string `json:"razorpay_signature"`
	PurchaseID        string `json:"purchase_id"`
}

func VerifyPayment(c *fiber.Ctx) error {
	var body verifyReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	payload := body.RazorpayOrderID + "|" + body.RazorpayPaymentID
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	expected := hex.EncodeToString(h.Sum(nil))
	if expected != body.RazorpaySignature {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid signature"})
	}
	pid, err := uuid.Parse(body.PurchaseID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid purchase id"})
	}
	var p models.Purchase
	if err := store.DB.Where("id = ?", pid).First(&p).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "purchase not found"})
	}
	p.PaymentID = body.RazorpayPaymentID
	p.Status = "paid"
	store.DB.Save(&p)
	return c.JSON(fiber.Map{"ok": true})
}
