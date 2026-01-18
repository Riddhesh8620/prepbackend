package store

import (
	"log"
	"prepbackend/internal/models"
)

func RunMigrations() error {
	log.Println("Running database migrations...")
	_ = DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	_ = DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")

	return DB.AutoMigrate(
		&models.User{},
		&models.Course{},
		&models.Topic{},
		&models.Purchase{},
	)
}
