package routes

import (
	"smart-pantry/backend/internal/controllers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupAPIRoutes(e *echo.Echo, pantryController *controllers.PantryController) {
	api := e.Group("/api")
	api.Use(middleware.Logger())
	api.Use(middleware.Recover())

	pantry := api.Group("/pantry")
	pantry.GET("", pantryController.GetPantryItems)
	pantry.POST("", pantryController.CreatePantryItem)
	pantry.PUT("/:id", pantryController.UpdatePantryItem)
	pantry.DELETE("/:id", pantryController.DeletePantryItem)
}
