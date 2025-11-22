package handlers

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func UserDashboard(c *fiber.Ctx) error {
	uidStr := c.Locals("user_id").(string)
	uid, _ := uuid.Parse(uidStr)

	var purchases []models.Purchase
	store.DB.Where("user_id = ? AND status = ?", uid, "paid").Find(&purchases)

	type entry struct {
		PurchaseID string  `json:"purchase_id"`
		TopicID    *string `json:"topic_id"`
	}
	var out []entry
	for _, p := range purchases {
		var tid *string
		if p.TopicID != nil {
			s := p.TopicID.String()
			tid = &s
		}
		out = append(out, entry{PurchaseID: p.ID.String(), TopicID: tid})
	}
	return c.JSON(fiber.Map{
		"user_id":   uidStr,
		"purchases": out,
	})
}
