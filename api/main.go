package main

import (
	"log"
	"os"
	"url-shortner/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func setUpRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	// Try to load .env file, but don't panic if it doesn't exist (for Docker)
	godotenv.Load(".env")
	
	// Print environment variables for debugging
	log.Printf("Environment variables: DB_ADDR=%s, APP_PORT=%s, DOMAIN=%s\n", 
		os.Getenv("DB_ADDR"), os.Getenv("APP_PORT"), os.Getenv("DOMAIN"))
	
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Log the error
			log.Printf("Error: %v\n", err)
			// Return a 500 response
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		},
	})
	
	// Add middleware for logging requests
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Request: %s %s\n", c.Method(), c.Path())
		return c.Next()
	})
	
	setUpRoutes(app)
	
	// Log server start
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = ":3000" // Default port if not set
	}
	log.Printf("Server is starting on port %s\n", port)
	
	// Handle errors
	if err := app.Listen(port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
