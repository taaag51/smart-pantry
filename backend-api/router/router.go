package router

import (
	"log"
	"net/http"
	"os"

	"github.com/taaag51/smart-pantry/backend-api/controller"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, fc controller.IFoodItemController, rc controller.IRecipeController) *echo.Echo {
	e := echo.New()

	// CORSミドルウェアの設定
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", os.Getenv("FRONTEND_URL")},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-CSRF-Token"},
		ExposeHeaders:    []string{"X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	// セキュリティヘッダーの設定
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         31536000,
		HSTSPreloadEnabled: true,
	}))

	// CSRFミドルウェアの設定
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookiePath:     "/",
		CookieHTTPOnly: true,
		CookieSameSite: http.SameSiteLaxMode,
		TokenLookup:    "header:X-CSRF-Token",
	}))

	// 認証関連（認証不要なエンドポイント）
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
	e.POST("/refresh-token", uc.RefreshToken) // 新しいリフレッシュトークンエンドポイント

	// トークン検証エンドポイント（認証ミドルウェアなし）
	e.GET("/verify-token", uc.VerifyToken)

	// JWT認証が必要なルート
	api := e.Group("/api")
	api.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
		SigningKey:    []byte(os.Getenv("SECRET")),
		TokenLookup:   "header:Authorization",
		ContextKey:    "user",
		SigningMethod: "HS256",
		ErrorHandler: func(c echo.Context, err error) error {
			// JWTミドルウェアのエラーをログ出力
			log.Printf("JWT認証エラー: %v", err)
			log.Printf("リクエストヘッダー: %v", c.Request().Header)
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"status":  "error",
				"message": "未認証",
			})
		},
	}))

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
