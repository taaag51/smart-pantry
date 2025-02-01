package usecase

import (
	"backend-api/model"
	"backend-api/repository"
	"backend-api/services"
	"errors"
	"fmt"
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
		return "", errors.New("食材が登録されていません。食材を追加してからレシピを取得してください")
	}

	// レシピを生成
	recipe, err := ru.gs.GenerateRecipe(foodItems)
	if err != nil {
		// Geminiサービスのエラーをログに出力
		fmt.Printf("Geminiサービスエラー: %v\n", err)
		return "レシピの生成中にエラーが発生しました。しばらく待ってから再試行してください。", nil
	}

	if recipe == "" {
		return "レシピを生成できませんでした。別の食材を試してみてください。", nil
	}

	return recipe, nil
}
