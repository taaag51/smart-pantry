package main

import (
	"backend-api/controller"
	"backend-api/db"
	"backend-api/repository"
	"backend-api/router"
	"backend-api/services"
	"backend-api/usecase"
	"backend-api/validator"
	"log"
)

func main() {
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
