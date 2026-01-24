package handlers

import (
	"net/http"

	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type createCourseReq struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price_in_paisa"`
	CategoryID  uuid.UUID `json:"category_id"`
	Thumbnail   string    `json:"thumbnail"`
}

func GetCourses(c *fiber.Ctx) error {
	var courses []models.Course
	store.DB.Preload("Topics").Find(&courses)
	return c.JSON(courses)
}

func GetCourse(c *fiber.Ctx) error {
	id := c.Params("id")
	var course models.Course
	if err := store.DB.Preload("Topics").Where("id = ?", id).First(&course).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "course not found"})
	}
	return c.JSON(course)
}

func AdminCreateCourse(c *fiber.Ctx) error {
	var body createCourseReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	course := models.Course{
		Title:       body.Title,
		Description: body.Description,
		CategoryID:  body.CategoryID,
	}
	if err := store.DB.Create(&course).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}
	return c.JSON(course)
}

type createTopicReq struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	OrderIndex  int     `json:"order_index"`
	Price       float32 `json:"price"`
}

func AdminCreateTopic(c *fiber.Ctx) error {
	courseID := c.Params("id")
	uid, err := uuid.Parse(courseID)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid course id"})
	}

	var body createTopicReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	t := models.Topic{
		CourseID:    uid,
		Title:       body.Title,
		Description: body.Description,
		OrderIndex:  body.OrderIndex,
		Price:       body.Price,
	}

	if err := store.DB.Create(&t).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	return c.JSON(t)
}
