package main

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/taaag51/smart-pantry/backend-api/controller"
	"github.com/taaag51/smart-pantry/backend-api/db"
	"github.com/taaag51/smart-pantry/backend-api/repository"
	"github.com/taaag51/smart-pantry/backend-api/router"
	"github.com/taaag51/smart-pantry/backend-api/services"
	"github.com/taaag51/smart-pantry/backend-api/usecase"
	"github.com/taaag51/smart-pantry/backend-api/validator"
)

func main() {
	// プロジェクトルートの .env ファイルを Overload で読み込み（上書き）
	if err := godotenv.Overload(filepath.Join("..", ".env")); err != nil {
		log.Println("Warning: project root .env file not loaded:", err)
	}
	db := db.NewDB()
	userValidator := validator.NewUserValidator()

	// リポジトリの初期化
	userRepository := repository.NewUserRepository(db)
	foodItemRepository := repository.NewFoodItemRepository(db)

	// サービスの初期化
	geminiService, err := services.NewGeminiService()
	if err != nil {
		log.Fatalf("Failed to initialize Gemini service: %v", err)
	}

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	foodItemUsecase := usecase.NewFoodItemUsecase(foodItemRepository)
	recipeUsecase := usecase.NewRecipeUsecase(foodItemRepository, geminiService)

	// コントローラーの初期化
	userController := controller.NewUserController(userUsecase)
	foodItemController := controller.NewFoodItemController(foodItemUsecase)
	recipeController := controller.NewRecipeController(recipeUsecase)

	// ルーターの設定
	e := router.NewRouter(userController, foodItemController, recipeController)
	e.Logger.Fatal(e.Start(":8080"))
}
