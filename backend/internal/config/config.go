package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() error {
	godotenv.Load() // load .env (ignore error if not present)

	if os.Getenv("DATABASE_URL") == "" {
		return errors.New("missing env: DATABASE_URL")
	}
	if os.Getenv("JWT_SECRET") == "" {
		return errors.New("missing env: JWT_SECRET")
	}

	return nil
}
