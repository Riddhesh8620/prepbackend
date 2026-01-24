package handlers

import (
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/google/uuid"
)

type createTopicReq struct {
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
