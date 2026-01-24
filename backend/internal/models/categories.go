package models

type Category struct {
	Base
	Title       string `gorm:"size:150;not null;uniqueIndex;" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	IconName    string `gorm:"size:100" json:"icon_name"`
	Color       string `gorm:"size:50" json:"color"`
	// Relations
	Courses []Course `gorm:"foreignKey:CategoryID" json:"courses"`
}
