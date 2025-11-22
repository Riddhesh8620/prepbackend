package store

import "prepbackend/internal/models"

func RunMigrations() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Topic{},
		&models.Purchase{},
	)
}
