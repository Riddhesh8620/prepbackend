package handlers

import (
	"net/http"
	"os"

	"prepbackend/internal/models"
	"prepbackend/internal/store"
	"prepbackend/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type signupReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func SignUp(c *fiber.Ctx) error {
	var body signupReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	// duplicate check
	var ex models.User
	if err := store.DB.Where("email = ?", body.Email).First(&ex).Error; err == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "email exists"})
	}
	hash, _ := utils.HashPassword(body.Password)
	u := models.User{
		ID:           uuid.New(),
		Name:         body.Name,
		Email:        body.Email,
		PasswordHash: hash,
		Role:         "student",
	}
	if err := store.DB.Create(&u).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	token, _ := utils.CreateJWT(u.ID.String(), u.Role)
	return c.JSON(fiber.Map{"token": token})
}

func Login(c *fiber.Ctx) error {
	var body loginReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	var u models.User
	if err := store.DB.Where("email = ?", body.Email).First(&u).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := utils.CheckPasswordHash(u.PasswordHash, body.Password); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, _ := utils.CreateJWT(u.ID.String(), u.Role)
	return c.JSON(fiber.Map{"token": token})
}

// CreateDefaultAdminIfNotExists seeds admin if env provided
func CreateDefaultAdminIfNotExists() error {
	adminEmail := os.Getenv("DEFAULT_ADMIN_EMAIL")
	adminPass := os.Getenv("DEFAULT_ADMIN_PASSWORD")
	if adminEmail == "" || adminPass == "" {
		return nil
	}
	var u models.User
	if err := store.DB.Where("email = ?", adminEmail).First(&u).Error; err == nil {
		return nil // exists
	}
	hash, _ := utils.HashPassword(adminPass)
	a := models.User{
		ID:           uuid.New(),
		Name:         "Admin",
		Email:        adminEmail,
		PasswordHash: hash,
		Role:         "admin",
	}
	return store.DB.Create(&a).Error
}
