package usecase

import (
	"errors"
	"fmt"

	"github.com/taaag51/smart-pantry/backend-api/model"
	"github.com/taaag51/smart-pantry/backend-api/repository"
	"github.com/taaag51/smart-pantry/backend-api/services"
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
	// ユーザーIDでフィルタリングした食材のみを使用
	var userFoodItems []model.FoodItem
	for _, item := range foodItems {
		if item.UserId == userId {
			userFoodItems = append(userFoodItems, item)
		}
	}

	// フィルタリング後の食材が0個の場合
	if len(userFoodItems) == 0 {
		return "", errors.New("登録されている食材がありません")
	}

	recipe, err := ru.gs.GenerateRecipe(userFoodItems)
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
