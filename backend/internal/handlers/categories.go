package handlers

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
)

func GetCategory(c *fiber.Ctx) error {
	var Category []models.Category
	store.DB.Find(&Category)
	return c.JSON(Category)
}
