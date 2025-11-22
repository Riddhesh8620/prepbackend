package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"

	"prepbackend/internal/config"
	"prepbackend/internal/handlers"
	"prepbackend/internal/middleware"
	"prepbackend/internal/store"
)

func main() {
	// load config (validates required envs)
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("config load: %v", err)
	}

	// connect DB
	if err := store.ConnectDB(); err != nil {
		log.Fatalf("db connect: %v", err)
	}

	// run migrations
	if err := store.RunMigrations(); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// seed admin if env present
	if err := handlers.CreateDefaultAdminIfNotExists(); err != nil {
		log.Printf("seed admin: %v", err)
	}

	app := fiber.New(fiber.Config{
		AppName: "PrepBackend",
	})

	api := app.Group("/api")

	// auth
	api.Post("/auth/signup", handlers.SignUp)
	api.Post("/auth/login", handlers.Login)

	// public
	api.Get("/courses", handlers.GetCourses)
	api.Get("/courses/:id", handlers.GetCourse)

	// user
	user := api.Group("/user")
	user.Use(middleware.RequireAuth)
	user.Get("/dashboard", handlers.UserDashboard)
	user.Post("/purchase/topic/:id", handlers.CreateTopicPurchase)
	user.Post("/payment/verify", handlers.VerifyPayment)

	// admin
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAuth, middleware.RequireAdmin)
	admin.Post("/courses", handlers.AdminCreateCourse)
	admin.Post("/courses/:id/topics", handlers.AdminCreateTopic)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
