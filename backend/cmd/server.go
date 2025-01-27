package main

import (
	"log"
	"smart-pantry/backend/configs"
	"smart-pantry/backend/internal/controllers"
	"smart-pantry/backend/internal/models"
	"smart-pantry/backend/internal/routes"
	"smart-pantry/backend/internal/services"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs.LoadEnv()

	db, err := gorm.Open(sqlite.Open("smart_pantry.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// AutoMigrate
	db.AutoMigrate(&models.PantryItem{})

	pantryService := services.NewPantryService(db)
	pantryController := controllers.NewPantryController(pantryService)

	e := echo.New()
	routes.SetupAPIRoutes(e, pantryController)

	e.Logger.Fatal(e.Start(":8080"))
}
