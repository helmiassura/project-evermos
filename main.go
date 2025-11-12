package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"evermos-project/config"
	"evermos-project/routes"
	"evermos-project/utils"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found")
	}

	// Initialize database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate database
	if err := config.MigrateDB(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			// don't leak internals in production — here we keep the original behavior
			return c.Status(code).JSON(fiber.Map{
				"status":  false,
				"message": "Internal Server Error",
				"errors":  []string{err.Error()},
				"data":    nil,
			})
		},
	})

	// Middleware
	app.Use(cors.New())
	app.Use(logger.New())

	// Serve static files (uploads)
	app.Static("/uploads", "./uploads")

	// Root health-check
	app.Get("/", func(c *fiber.Ctx) error {
		return utils.RespondJSON(c, fiber.StatusOK, true, "API is running", nil, fiber.Map{
			"service": "Evermos API",
			"status":  "ok",
		})
	})

	// Setup routes
	routes.SetupRoutes(app, db)

	// Fallback 404 JSON (must be after routes)
	app.Use(func(c *fiber.Ctx) error {
		return utils.RespondJSON(c, fiber.StatusNotFound, false, "Not Found", []string{"Endpoint not found"}, nil)
	})

	// Get port from env
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start server
	log.Printf("Server running on port %s\n", port)
	log.Fatal(app.Listen(":" + port))
}
