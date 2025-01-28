package main

import (
	"log"
	"net/http"

	"smart-pantry/backend/configs"
	"smart-pantry/backend/internal/models"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database connection
	configs.ConnectDatabase()

	// Perform database migrations
	if err := configs.DB.AutoMigrate(&models.FoodItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, Smart Pantry!")
	})

	// Start server
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("Shutting down the server: %v", err)
	}
}
