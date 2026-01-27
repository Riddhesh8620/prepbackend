package handlers

import (
	"net/http"
	"prepbackend/internal/models"
	"prepbackend/internal/store"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type createTopicReq struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int8      `json:"orderIndex"`
	Price       float32   `json:"price"`
	Duration    string    `json:"duration"` // in minutes
	CourseID    uuid.UUID `json:"courseId"`
}

func AdminCreateTopic(body []createTopicReq, courseID *uuid.UUID) bool {

	var topics []models.Topic
	for _, item := range body {
		t := models.Topic{
			CourseID:    *courseID,
			Title:       item.Title,
			Description: item.Description,
			OrderIndex:  item.OrderIndex,
			Price:       item.Price,
			Duration:    item.Duration,
		}
		topics = append(topics, t)
	}

	if err := store.DB.Create(&topics).Error; err != nil {
		return false
	}
	return true
}

func AdminUpdateTopics(txContext *gorm.DB, topics []createTopicReq) bool {
	var err error
	for _, t := range topics {
		// Skip empty entries
		if t.Title == "" {
			continue
		}

		topic := models.Topic{
			Title:    t.Title,
			Price:    t.Price,
			Duration: t.Duration,
			CourseID: t.CourseID,
		}
		id, parseErr := uuid.Parse(t.ID)

		if parseErr == nil && id != uuid.Nil {
			// 1. UPDATE: Topic has an ID, so update the existing record
			// Use Omit("id") to ensure the primary key isn't accidentally changed
			err = txContext.Model(&models.Topic{}).
				Where("id = ? AND course_id = ?", id, t.CourseID).
				Updates(topic).Error

			if err != nil {
				return false
			}
		} else {
			// 2. CREATE: No valid ID present, generate new record
			err = txContext.Create(&topic).Error
			if err != nil {
				return false
			}
		}
	}
	return err == nil
}

func AdminUpdateTopicInternal(c *fiber.Ctx) error {
	var dbError error
	tx := store.DB.Begin()

	price, _ := strconv.ParseFloat(c.FormValue("price"), 32)
	duration := c.FormValue("duration")
	courseId, _ := uuid.Parse(c.FormValue("courseId"))
	topicId := c.FormValue("id")
	title := c.FormValue("title")

	model := models.Topic{
		Title:       title,
		Description: "",
		Price:       float32(price),
		CourseID:    courseId,
		Duration:    duration,
	}

	if topicId != "" && title != "" {
		id, _ := uuid.Parse(topicId)

		dbError = tx.Model(&models.Topic{}).
			Where("id = ? AND course_id = ?", id, courseId).
			Updates(model).Error
	} else {
		dbError = tx.Create(&model).Error
	}

	dbError = tx.Commit().Error

	if dbError != nil {
		tx.Rollback()
	}

	return c.Status(http.StatusOK).JSON(model)
}
