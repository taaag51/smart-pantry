package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(tc controller.ITaskController, uc controller.IUserController, fc controller.IFoodItemController, rc controller.IRecipeController) *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定を修正
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteDefaultMode,
		TokenLookup:    "header:X-CSRF-Token",
	}))

	// 認証関連
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)

	// JWT認証が必要なルート
	api := e.Group("")
	api.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
		SigningKey:  []byte(os.Getenv("SECRET")),
		TokenLookup: "cookie:token",
	}))

	// トークン検証エンドポイント
	api.GET("/verify-token", uc.VerifyToken)

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
