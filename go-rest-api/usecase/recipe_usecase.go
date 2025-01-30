package usecase

import (
	"fmt"
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
		return "", fmt.Errorf("食材の取得に失敗しました: %v", err)
	}

	// 食材が存在しない場合はエラーを返す
	if len(foodItems) == 0 {
		return "", fmt.Errorf("食材が登録されていません")
	}

	// レシピを生成
	recipe, err := ru.gs.GenerateRecipe(foodItems)
	if err != nil {
		return "", fmt.Errorf("レシピの生成に失敗しました: %v", err)
	}

	return recipe, nil
}
