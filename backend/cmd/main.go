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
	// if err := store.RunMigrations(); err != nil {
	// 	log.Fatalf("migrate: %v", err)
	// }

	// seed admin if env present
	// if err := handlers.CreateDefaultAdminIfNotExists(); err != nil {
	// 	log.Printf("seeded admin: %v", err)
	// }

	// seed default categories
	// if err := handlers.AdminCreateDefaultCategory(); err != nil {
	// 	log.Printf("seeded categories: %v", err)
	// }

	app := fiber.New(fiber.Config{
		AppName: "PrepBackend",
		Prefork: false,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:8080, https://kodingkraze6-dashboard-903239c8.vercel.app", // Add your Vercel URL
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	api := app.Group("/")
	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("PrepBackend is running")
	})

	api = app.Group("/api")

	// auth
	api.Post("/auth/signup", handlers.SignUp)
	api.Post("/auth/login", handlers.Login)
	api.Post("/auth/logout", handlers.Logout)
	api.Post("auth/send/otp", handlers.SendOTPEmail)
	api.Post("auth/otp/verify", handlers.VerifyOTP)
	// public
	courseViewGrp := api.Group("/courses")
	courseViewGrp.Get("/", handlers.GetCourses)
	courseViewGrp.Get("/get-by-id/:id", handlers.GetCourse)
	courseViewGrp.Get("/:categoryId", handlers.GetCoursesByCategory)

	// Category
	api.Get("/categories", handlers.GetCategory)
	api.Get("/categories/:id", handlers.GetCategoryById)
	categoryCreateGroup := api.Group("/categories")
	categoryCreateGroup.Use(middleware.RequireAuth, middleware.RequireAdmin)
	categoryCreateGroup.Post("/save", handlers.SaveCategory)

	// user
	user := api.Group("/user")
	user.Use(middleware.RequireAuth)
	user.Get("/dashboard", handlers.UserDashboard)
	user.Post("/purchase/topic/:id", handlers.CreateTopicPurchase)
	user.Post("/payment/verify", handlers.VerifyPayment)

	// admin

	admin := api.Group("/courses")
	admin.Use(middleware.RequireAuth, middleware.RequireAdmin)
	admin.Post("/save", handlers.CreateCourse)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
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
