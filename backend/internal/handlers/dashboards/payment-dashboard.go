package dashboards

import (
	"net/http"
	"net/smtp"
	"os"
	"prepbackend/cmd/email"
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

func SendPaymentSuccesfullEmail(c *fiber.Ctx, sessionID *uuid.UUID) error {
	GMAIL_ID := os.Getenv("GMAIL_ID")
	APP_PASSWORD := os.Getenv("APP_PASSWORD")
	DASHBOARD_URL := os.Getenv("DASHBOARD_URL")

	var session models.UserPaymentSession
	if err := store.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "session not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	var user models.User
	if err := store.DB.First(&user, "id = ?", session.UserID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	var topicsGrp []models.TopicGroup
	dbError := store.DB.Raw(`select c.id, ti.id, ti.user_payment_session_id, t.title as topic_title, c.title as course_title from topic_inventories as ti 
	inner join topics as t on t.id::uuid = ti.topic_id::uuid
	inner join courses as c on c.id::uuid = t.course_id::uuid
	WHERE ti.user_payment_session_id = ?
	GROUP BY c.id::uuid ,ti.id::uuid, ti.user_payment_session_id::uuid, t.title, c.title;`,
		sessionID).
		Scan(&topicsGrp).
		Error

	if dbError != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	var courseProjectList []models.CourseProject
	dbError = store.DB.Raw(`SELECT c.title as course_title, ci.id, ci.user_payment_session_id from course_inventories as ci 
	inner join courses as c ON c.id::uuid = ci.course_id::uuid
	where ci.user_payment_session_id::uuid = ? 
	ORDER BY c.title ASC;`, sessionID).
		Scan(&courseProjectList).
		Error

	if dbError != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	groups := make(map[string]*models.CourseGroup)
	for _, row := range topicsGrp {
		if _, exists := groups[row.CourseTitle]; !exists {
			groups[row.CourseTitle] = &models.CourseGroup{
				Topics: []string{},
			}
		}
		groups[row.CourseTitle].Topics = append(groups[row.CourseTitle].Topics, row.TopicTitle)
	}

	var emailTopicList []email.TopicItem
	var emailCourseList []email.CourseItem

	for k, courseGrp := range groups {
		emailTopicList = append(emailTopicList, email.TopicItem{
			ParentCourse: k,
			Topics:       courseGrp.Topics,
		})
	}

	for _, c := range courseProjectList {
		emailCourseList = append(emailCourseList, email.CourseItem{
			CourseTitle: c.CourseTitle,
		})
	}

	paymentSuccessPayload := email.ActivationPayload{
		Topics:       emailTopicList,
		DashboardURL: DASHBOARD_URL,
		Courses:      emailCourseList,
	}

	htmlBody := email.BuildActivationEmail(paymentSuccessPayload)
	from := GMAIL_ID
	pass := APP_PASSWORD
	subject := "Subject: 🔑 Your Order Is Confirmed | Access Unlocked\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := []byte(subject + mime + htmlBody)
	return smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{user.Email}, []byte(msg))
}

// Helper to prevent duplicate topics in the list
func uniqueStrings(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
