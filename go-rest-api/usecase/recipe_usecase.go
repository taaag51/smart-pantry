package usecase

import (
	"go-rest-api/model"
	"go-rest-api/repository"
	"go-rest-api/services"
)

type IRecipeUsecase interface {
	GetRecipeSuggestions(userId uint) (string, error)
}

type recipeUsecase struct {
	fr repository.IFoodItemRepository
	gs services.IGeminiService
}

func NewRecipeUsecase(fr repository.IFoodItemRepository, gs services.IGeminiService) IRecipeUsecase {
	return &recipeUsecase{fr, gs}
}

func (ru *recipeUsecase) GetRecipeSuggestions(userId uint) (string, error) {
	// ユーザーの食材一覧を取得
	var foodItems []model.FoodItem
	if err := ru.fr.GetAllFoodItems(&foodItems); err != nil {
		return "", err
	}

	// レシピを生成
	recipe, err := ru.gs.GenerateRecipe(foodItems)
	if err != nil {
		return "", err
	}

	return recipe, nil
}
