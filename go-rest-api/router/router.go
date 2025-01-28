package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(tc controller.ITaskController, uc controller.IUserController, fc controller.IFoodItemController, rc controller.IRecipeController) *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept,
			echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE"},
		AllowCredentials: true,
	}))
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieDomain:   os.Getenv("API_DOMAIN"),
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteNoneMode,
	}))

	// 認証関連
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	// JWT認証が必要なルート
	api := e.Group("")
	api.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))

	// タスク関連
	tasks := api.Group("/tasks")
	tasks.GET("", tc.GetAllTasks)
	tasks.GET("/:taskId", tc.GetTaskById)
	tasks.POST("", tc.CreateTask)
	tasks.PUT("/:taskId", tc.UpdateTask)
	tasks.DELETE("/:taskId", tc.DeleteTask)

	// 食材関連
	foodItems := api.Group("/food-items")
	foodItems.GET("", fc.GetAllFoodItems)
	foodItems.GET("/:id", fc.GetFoodItemById)
	foodItems.POST("", fc.CreateFoodItem)
	foodItems.PUT("/:id", fc.UpdateFoodItem)
	foodItems.DELETE("/:id", fc.DeleteFoodItem)

	// レシピ関連
	recipes := api.Group("/recipes")
	recipes.GET("/suggestions", rc.GetRecipeSuggestions)

	return e
}
