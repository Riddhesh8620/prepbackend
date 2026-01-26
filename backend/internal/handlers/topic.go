package handlers

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type createTopicReq struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	OrderIndex  int8      `json:"order_index"`
	Price       float32   `json:"price"`
	Duration    string    `json:"duration"` // in minutes
	CourseID    uuid.UUID `json:"course_id"`
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

func AdminUpdateTopic(txContext *gorm.DB, topics []createTopicReq) bool {
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

		if t.ID != uuid.Nil {
			// 1. UPDATE: Topic has an ID, so update the existing record
			// Use Omit("id") to ensure the primary key isn't accidentally changed
			err = txContext.Model(&models.Topic{}).
				Where("id = ? AND course_id = ?", t.ID, t.CourseID).
				Updates(topic).Error

			if err != nil {
				return false
			}
		} else {
			// 2. CREATE: No ID present, generate new record
			if err = txContext.Create(&topic).Error; err != nil {
				err = nil
			}
		}
	}

	if err == nil {
		return true
	} else {
		return false
	}
}
