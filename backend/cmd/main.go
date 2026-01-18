package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"prepbackend/internal/config"
	"prepbackend/internal/handlers"
	"prepbackend/internal/middleware"
	"prepbackend/internal/store"
)

func main() {

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting wd:", err)
	} else {
		log.Println("PWD =", wd)
	}

	// FORCE load .env from root
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// load config (validates required envs)
	if err := config.LoadConfig(); err != nil {
		log.Println(".env not found, relying on system env")
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
		Prefork: false,
	})

	api := app.Group("/")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("PrepBackend is running")
	})

	api = app.Group("/api")

	// auth
	api.Post("/auth/signup", handlers.SignUp)
	api.Post("/auth/login", handlers.Login)
	api.Post("/auth/logout", handlers.Logout)

	// public
	api.Get("/courses", handlers.GetCourses)
	api.Get("/courses/:id", handlers.GetCourse)

	// Category
	api.Get("/Category", handlers.GetCategory)

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

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
}
