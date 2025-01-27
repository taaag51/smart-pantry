package routes

import (
	"smart-pantry/backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			authController := controllers.NewAuthController()
			auth.POST("/login", authController.Login)
			auth.POST("/register", authController.Register)
			auth.POST("/logout", authController.Logout)
		}

		pantry := api.Group("/pantry")
		{
			pantry.GET("", controllers.GetPantryItems)
			pantry.POST("", controllers.AddPantryItem)
			pantry.PUT("/:id", controllers.UpdatePantryItem)
			pantry.DELETE("/:id", controllers.DeletePantryItem)
		}
	}
}
