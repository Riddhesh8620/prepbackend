package models

type User struct {
	Base
	Name         string `gorm:"size:255;not null"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:50;default:'student'"`
}
