package store

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	// attempt to create extension if superuser; ignore error
	_ = DB.Exec("CREATE EXTENSION IF NOT EXISTS pgcrypto;")
	return nil
}
