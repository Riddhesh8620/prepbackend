package dashboards

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
)

func GetAllPaymentSessions(c *fiber.Ctx) error {
	var dtoList []models.GetPaymentDashboardDto

	// 1. We query the sessions table and join the users table
	// 2. We select only the fields required for the DTO
	// 3. We use the 'sessions' alias and 'users' alias for clarity
	err := store.DB.Table("user_payment_sessions as s").
		Select("s.id, u.email as user_email, u.name as name, s.payable_amount as amount, s.status as status").
		Joins("left join users u on u.id = s.user_id").
		Order("s.created_on desc").
		Scan(&dtoList).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch payment dashboard data",
		})
	}

	// 4. Return the DTO list instead of the raw model
	return c.Status(fiber.StatusOK).JSON(dtoList)
}
