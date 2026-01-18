package models

type Category struct {
	Base
	Title       string `gorm:"size:150;not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Relations
	Courses []Course `gorm:"foreignKey:CategoryID" json:"courses"`
}
