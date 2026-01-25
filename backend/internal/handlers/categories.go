package handlers

import (
	"net/http"
	"prepbackend/internal/models"
	"prepbackend/internal/store"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type CategorySeed struct {
	Color    string
	IconName string
}

type CategoryDto struct {
	ID          string `json:"id"`
	Title       string `json:"name"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	IconName    string `json:"icon"`
	Color       string `json:"color"`
	CourseCount int16  `json:"courseCount"`
}

type CreateCateogoryReq struct {
	Title       string `json:"name"`
	Description string `json:"description" nullable:"true"`
	IconName    string `json:"icon"`
	ID          string `json:"id" nullable:"true"`
	Color       string `json:"color"`
}

func GetCategoryByIdInternal(categoryId string) (models.Category, error) {
	var dto models.Category
	id, err := uuid.Parse(categoryId)
	if err != nil || id == uuid.Nil {
		return dto, err
	}
	// Use .First() to retrieve the category data directly into the DTO.
	// This will map the fields from the categories table to your struct.
	err = store.DB.
		Model(&models.Category{}).
		First(&dto, "id = ?", id).Error

	if err != nil {
		// Return the empty DTO and the error to the calling function
		return dto, err
	}

	return dto, nil
}

func GetCategoryById(c *fiber.Ctx) error {
	// 1. Parse the ID from parameters
	categoryId, err := uuid.Parse(c.Params("id"))
	if err != nil || categoryId == uuid.Nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid or missing Category ID",
		})
	}

	// 2. Define a single DTO (not a slice)
	var dto CategoryDto

	// 3. Query with a Where clause and a subquery for the count
	// This allows you to get all category info + course count in one go
	err = store.DB.Model(&models.Category{}).
		Select("categories.*, (SELECT COUNT(*) FROM courses WHERE courses.category_id = categories.id) as course_count").
		Where("id = ?", categoryId).
		First(&dto).Error

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Category not found",
		})
	}
	return c.JSON(dto)
}

func GetCategory(c *fiber.Ctx) error {
	var dtos []CategoryDto

	// This query gets categories and counts courses in one single trip to the DB
	err := store.DB.Model(&models.Category{}).
		Select("categories.*, (SELECT COUNT(*) FROM courses WHERE courses.category_id = categories.id) as course_count").
		Where("categories.is_active = ?", true). // Added explicit table reference
		Find(&dtos).Error

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch categories"})
	}

	return c.JSON(dtos)
}

func SaveCategory(c *fiber.Ctx) error {
	var body CreateCateogoryReq
	if err := c.BodyParser(&body); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	category := models.Category{
		Title:       body.Title,
		Description: body.Description,
		IconName:    body.IconName,
		IsActive:    true,
		Color:       body.Color,
	}

	if err := store.DB.Create(&category).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "db error"})
	}

	return c.JSON(body)
}

func AdminCreateDefaultCategory() error {
	var categories = map[string]CategorySeed{
		"Logical Reasoning":     {IconName: "Brain", Color: "hsl(262, 83%, 58%)"},
		"Analytical Reasoning":  {IconName: "BarChart3", Color: "hsl(199, 89%, 48%)"},
		"Verbal Ability":        {IconName: "MessageSquare", Color: "hsl(340, 82%, 52%)"},
		"Quantitative Aptitude": {IconName: "Calculator", Color: "hsl(142, 71%, 45%)"},
		"Web Development":       {IconName: "Code2", Color: "hsl(217, 91%, 60%)"},
		"DSA":                   {IconName: "Binary", Color: "hsl(25, 95%, 53%)"},
		"Data Science":          {IconName: "Database", Color: "hsl(280, 70%, 50%)"},
		"Cloud Computing":       {IconName: "Cloud", Color: "hsl(199, 89%, 48%)"},
		"Cybersecurity":         {IconName: "Shield", Color: "hsl(0, 70%, 50%)"},
		"Database":              {IconName: "Server", Color: "hsl(45, 90%, 50%)"},
		"Mobile Apps":           {IconName: "Smartphone", Color: "hsl(190, 80%, 45%)"},
		"DevOps":                {IconName: "GitBranch", Color: "hsl(170, 70%, 40%)"},
	}

	var Entities []models.Category
	var dbTitles []string

	err := store.DB.Model(&models.Category{}).Pluck("title", &dbTitles).Error
	if err != nil {
		return err
	}
	// 2. Convert slice to set for O(1) lookup
	categorySet := make(map[string]struct{})
	for _, title := range dbTitles {
		categorySet[title] = struct{}{}
	}

	for title, info := range categories {
		if _, found := categorySet[title]; !found {
			Entities = append(Entities, models.Category{
				Title:    title,
				IconName: info.IconName,
				Color:    info.Color, // Storing the HSL string
				IsActive: true,
			})
		}
	}

	if len(Entities) > 0 {
		return store.DB.Create(&Entities).Error
	}
	return nil
}
