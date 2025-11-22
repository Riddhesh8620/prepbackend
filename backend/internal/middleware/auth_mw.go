package middleware

import (
	"strings"

	"prepbackend/internal/models"
	"prepbackend/internal/store"
	"prepbackend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func RequireAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid auth header"})
	}
	tok, err := utils.ParseJWT(parts[1])
	if err != nil || !tok.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid claims"})
	}
	sub, _ := claims["sub"].(string)
	role, _ := claims["role"].(string)

	// verify user exists
	var u models.User
	if err := store.DB.First(&u, "id = ?", sub).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	c.Locals("user_id", sub)
	c.Locals("role", role)
	return c.Next()
}

func RequireAdmin(c *fiber.Ctx) error {
	role := c.Locals("role")
	if role == nil || role.(string) != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "admin only"})
	}
	return c.Next()
}
