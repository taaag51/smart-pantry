package routes

import (
	"smart-pantry/backend/internal/controllers"
	"smart-pantry/backend/internal/services"

	"github.com/labstack/echo/v4"
)

// SetupAPIRoutes defines the API routes for the application.
func SetupAPIRoutes(e *echo.Echo, as *services.AuthService, pc *controllers.AuthController, ps *services.PantryService) {
	apiGroup := e.Group("/api")

	// Auth routes
	authController := controllers.NewAuthController(as)
	apiGroup.POST("/signup", authController.SignUp)
	apiGroup.POST("/login", authController.Login)

	// Pantry routes
	pantryController := controllers.NewPantryController(ps)
	pantryGroup := apiGroup.Group("/pantry")
	pantryGroup.GET("", pantryController.GetPantryItems)
	pantryGroup.GET("/:id", pantryController.GetPantryItem)
	pantryGroup.POST("", pantryController.AddPantryItem)
	pantryGroup.PUT("/:id", pantryController.UpdatePantryItem)
	pantryGroup.DELETE("/:id", pantryController.DeletePantryItem)
	+pantryGroup.GET("/recipe", pantryController.SuggestRecipe) // レシピ提案エンドポイント
}
