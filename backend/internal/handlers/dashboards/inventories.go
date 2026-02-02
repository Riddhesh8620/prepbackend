package dashboards

import (
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetUserInventory(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Locals("user_id").(string)) // From your Auth Middleware

	// Define structures to hold the joined data
	type CourseResponse struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Thumbnail string    `json:"thumbnail"`
		Slug      string    `json:"slug"`
	}

	type TopicResponse struct {
		ID          uuid.UUID `json:"id"`
		Title       string    `json:"title"`
		ParentTitle string    `json:"parent_title"` // From the Course table
		Slug        string    `json:"slug"`
	}

	var courses []CourseResponse
	var topics []TopicResponse

	// 1. Fetch Active Courses with Meta Data
	err = store.DB.Table("course_inventories as ci").
		Select("c.id, c.title, c.thumbnail").
		Joins("inner join courses c on c.id::uuid = ci.course_id::uuid").
		Where("ci.user_id = ? AND ci.is_active = ?", userID, true).
		Scan(&courses).Error

	// 2. Fetch Active Topics with Meta Data (including the course it belongs to)
	err = store.DB.Table("topic_inventories as ti").
		Select("t.id, t.title, c.title as parent_title").
		Joins("inner join topics t on t.id::uuid = ti.topic_id::uuid"). // Force both to UUID
		Joins("inner join courses c on c.id::uuid = t.course_id::uuid").
		Where("ti.user_id = ? AND ti.is_active = ?", userID, true).
		Scan(&topics).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error fetching library"})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
		"topics":  topics,
	})
}
