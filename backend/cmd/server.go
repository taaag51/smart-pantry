package main

import (
	"fmt"
	"os"

	"smart-pantry/backend/db"
	"smart-pantry/backend/internal/controllers"
	"smart-pantry/backend/internal/routes"
	"smart-pantry/backend/internal/services"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// Load .env file
	godotenv.Load()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Database initialization
	database := db.InitDB()
	pantryService := services.NewPantryService(database)
	authService := services.NewAuthService(database)
	authController := controllers.NewAuthController(authService)
	pantryController := controllers.NewPantryController(pantryService)

	// Routes
	routes.SetupAPIRoutes(e, authService, authController, pantryService)

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080" // Default port
	}
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
