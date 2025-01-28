package main

import (
	"go-rest-api/controller"
	"go-rest-api/db"
	"go-rest-api/repository"
	"go-rest-api/router"
	"go-rest-api/services"
	"go-rest-api/usecase"
	"go-rest-api/validator"
	"log"
)

func main() {
	db := db.NewDB()
	userValidator := validator.NewUserValidator()
	taskValidator := validator.NewTaskValidator()

	// リポジトリの初期化
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	foodItemRepository := repository.NewFoodItemRepository(db)

	// サービスの初期化
	geminiService, err := services.NewGeminiService()
	if err != nil {
		log.Fatalf("Failed to initialize Gemini service: %v", err)
	}

	// ユースケースの初期化
	userUsecase := usecase.NewUserUsecase(userRepository, userValidator)
	taskUsecase := usecase.NewTaskUsecase(taskRepository, taskValidator)
	foodItemUsecase := usecase.NewFoodItemUsecase(foodItemRepository)
	recipeUsecase := usecase.NewRecipeUsecase(foodItemRepository, geminiService)

	// コントローラーの初期化
	userController := controller.NewUserController(userUsecase)
	taskController := controller.NewTaskController(taskUsecase)
	foodItemController := controller.NewFoodItemController(foodItemUsecase)
	recipeController := controller.NewRecipeController(recipeUsecase)

	// ルーターの設定
	e := router.NewRouter(taskController, userController, foodItemController, recipeController)
	e.Logger.Fatal(e.Start(":8080"))
}
