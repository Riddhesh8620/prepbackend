package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateTopicPurchase(c *fiber.Ctx) error {
	topicID := c.Params("id")
	tuid, err := uuid.Parse(topicID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid topic id"})
	}
	var topic models.Topic
	if err := store.DB.Where("id = ?", tuid).First(&topic).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "topic not found"})
	}

	uidStr := c.Locals("user_id").(string)
	userUUID, _ := uuid.Parse(uidStr)

	amount := topic.Price
	if amount <= 0 {
		// free topic: create paid purchase
		p := models.Purchase{
			UserID:  userUUID,
			TopicID: &tuid,
			Amount:  0,
			Status:  "paid",
		}
		_ = store.DB.Create(&p)
		return c.JSON(fiber.Map{"free": true, "purchase_id": p.ID})
	}

	orderID, err := createRazorpayOrder(amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create order"})
	}

	p := models.Purchase{
		UserID:    userUUID,
		TopicID:   &tuid,
		Amount:    amount,
		PaymentID: orderID, // temporarily store order id
		Status:    "pending",
	}
	if err := store.DB.Create(&p).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	return c.JSON(fiber.Map{
		"order_id":    orderID,
		"amount":      amount,
		"key_id":      os.Getenv("RAZORPAY_KEY_ID"),
		"purchase_id": p.ID.String(),
	})
}

func createRazorpayOrder(amount float32) (string, error) {
	data := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  fmt.Sprintf("rcpt_%d", amount),
	}
	b, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", "https://api.razorpay.com/v1/orders", bytes.NewReader(b))
	req.SetBasicAuth(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var out map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if id, ok := out["id"].(string); ok {
		return id, nil
	}
	return "", fmt.Errorf("no order id in response")
}
