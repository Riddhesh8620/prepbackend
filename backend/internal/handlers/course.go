package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type courseDetail struct {
	ID            string         `json:"id"`
	Title         string         `json:"title"`
	Price         float32        `json:"price"`
	OriginalPrice float32        `json:"originalPrice"`
	Duration      string         `json:"duration"`
	Level         string         `json:"level"`
	Image         string         `json:"image"`
	Description   string         `json:"description"`
	Curriculum    []models.Topic `json:"curriculum"`
	CategoryName  string         `json:"categoryName"`
	CategoryColor string         `json:"categoryColor"`
	CategoryID    string         `json:"categoryId"`
	Topics        int64          `json:"topics"`
}

type createCourseReq struct {
	Title         string           `json:"title"`
	Description   string           `json:"description"`
	OriginalPrice float32          `json:"original_price"`
	Price         float32          `json:"price"`
	CategoryID    string           `json:"category_id"`
	Topics        []createTopicReq `json:"topics"`
	Level         string           `json:"level"`
	Duration      string           `json:"duration"`
}

type updateCourseReq struct {
	createCourseReq
	CourseId string `json:"courseId"`
}

type courseDto struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Thumbnail     string  `json:"image"`
	Price         float32 `json:"price"`
	OriginalPrice float32 `json:"originalPrice"`
	Level         string  `json:"level"`
	Duration      string  `json:"duration"`
	CategoryID    string  `json:"category_id"`
	Topics        int64   `json:"topics"`
	CategoryName  string  `json:"categoryName"`
	CategoryColor string  `json:"categoryColor"`
}

func GetCourses(c *fiber.Ctx) error {
	var courses []models.Course

	var dtoList []courseDetail

	if err := store.DB.
		Preload("Topics").
		Find(&courses).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "courses not found"})

	}
	for _, model := range courses {

		categoryDto, _ := GetCategoryByIdInternal(model.CategoryID.String())

		var topicsCount int64 = 0
		for _, t := range model.Topics {
			// Assuming your 'createTopicReq' uses string IDs or you're checking for non-empty titles
			if t.IsActive {
				topicsCount = topicsCount + 1
			}
		}

		dtoList = append(dtoList, courseDetail{
			ID:            model.ID.String(),
			Title:         model.Title,
			Price:         model.Price,
			OriginalPrice: model.OriginalPrice,
			Duration:      model.Duration,
			Level:         model.Level,
			Image:         model.Thumbnail,
			Description:   model.Description,
			CategoryName:  categoryDto.Title,
			CategoryColor: categoryDto.Color,
			CategoryID:    categoryDto.ID.String(),
			Topics:        topicsCount,
			Curriculum:    model.Topics,
		})
	}
	return c.JSON(dtoList)
}

func GetCourse(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))

	if err != nil || id == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid category id"})
	}

	var course models.Course
	if err := store.DB.Preload("Topics").Where("id = ?", id).First(&course).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "course not found"})
	}

	var courseDetail courseDetail
	courseDetail.Curriculum = course.Topics
	courseDetail.ID = course.ID.String()
	courseDetail.Description = course.Description
	courseDetail.Duration = course.Duration
	courseDetail.Image = course.Thumbnail
	courseDetail.Level = course.Level
	courseDetail.Price = course.Price
	courseDetail.OriginalPrice = course.OriginalPrice
	courseDetail.Title = course.Title

	return c.JSON(courseDetail)
}

func CreateCourse(c *fiber.Ctx) error {
	file, err := c.FormFile("image")
	var base64Image string
	if err == nil {
		openedFile, _ := file.Open()
		defer openedFile.Close()
		fileBytes := make([]byte, file.Size)
		openedFile.Read(fileBytes)
		mimeType := file.Header.Get("Content-Type")
		base64Image = fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(fileBytes))
	}

	// 2. Parse Topics JSON String
	var topics []createTopicReq
	if err := json.Unmarshal([]byte(c.FormValue("topics")), &topics); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid topics format"})
	}

	// 3. Convert Numeric strings
	price, _ := strconv.ParseFloat(c.FormValue("price"), 32)
	origPrice, _ := strconv.ParseFloat(c.FormValue("original_price"), 32)
	catID, _ := uuid.Parse(c.FormValue("category_id"))
	level := c.FormValue("level")
	duration := c.FormValue("duration")
	// 4. Build Model
	course := models.Course{
		Title:         c.FormValue("title"),
		Description:   c.FormValue("description"),
		Price:         float32(price),
		OriginalPrice: float32(origPrice),
		CategoryID:    catID,
		Thumbnail:     base64Image,
		IsActive:      true,
		Level:         level,
		Duration:      duration,
	}

	// 5. Save Course & Topics
	if err := store.DB.Create(&course).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "DB Save Failed"})
	}

	AdminCreateTopic(topics, &course.ID)

	return c.Status(201).JSON(course)
}

func GetCourseIDsByCategory(categoryID *uuid.UUID) (int16, error) {
	var count int64 // GORM Count requires int64

	// This executes: SELECT count(*) FROM courses WHERE category_id = '...'
	err := store.DB.Model(&models.Course{}).
		Where("category_id = ?", categoryID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	// Convert int64 to int16 for your specific return type
	return int16(count), nil
}

func GetCoursesByCategory(c *fiber.Ctx) error {
	var courseDtoList []courseDto

	categoryID, err := uuid.Parse(c.Params("categoryId"))
	if err != nil || categoryID == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid category id"})
	}

	err = store.DB.Model(&models.Course{}).
		Select(`
            courses.*, 
            (SELECT COUNT(*) FROM topics 
             WHERE topics.course_id = courses.id 
             AND topics.is_active = true) as topics
        `).
		Where("category_id = ?", categoryID).
		Scan(&courseDtoList).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not fetch courses",
		})
	}
	return c.JSON(courseDtoList)
}

func UpdateCourse(c *fiber.Ctx) error {

	file, err := c.FormFile("image")
	var base64Image string
	if err == nil {
		openedFile, _ := file.Open()
		defer openedFile.Close()
		fileBytes := make([]byte, file.Size)
		openedFile.Read(fileBytes)
		mimeType := file.Header.Get("Content-Type")
		base64Image = fmt.Sprintf("data:%s;base64,%s", mimeType, base64.StdEncoding.EncodeToString(fileBytes))
	}

	// Start a transaction to ensure atomicity
	tx := store.DB.Begin()
	courseID, err := uuid.Parse(c.FormValue("id"))
	// 1. Update Course Metadata
	var course models.Course
	if err := tx.First(&course, "id = ?", courseID).Error; err != nil {
		tx.Rollback()
		return c.Status(404).JSON(fiber.Map{"error": "Course ID not found "})
	}

	price, _ := strconv.ParseFloat(c.FormValue("price"), 32)
	origPrice, _ := strconv.ParseFloat(c.FormValue("original_price"), 32)
	catID, _ := uuid.Parse(c.FormValue("category_id"))
	level := c.FormValue("level")
	duration := c.FormValue("duration")

	course.Title = c.FormValue("title")
	course.Description = c.FormValue("description")
	course.Price = float32(price)
	course.OriginalPrice = float32(origPrice)
	course.CategoryID = catID
	course.Level = level
	course.Duration = duration

	if base64Image != "" {
		course.Thumbnail = base64Image
	}

	if err := tx.Save(&course).Error; err != nil {
		tx.Rollback()
		return c.Status(500).JSON(fiber.Map{"error": "failed to update course"})
	}

	// Commit the transaction
	tx.Commit()

	return c.Status(200).JSON(fiber.Map{
		"message": "Course and curriculum updated successfully",
		"id":      courseID,
	})
}
