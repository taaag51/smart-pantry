package services

import (
	"go-rest-api/model"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRecipe(t *testing.T) {
	// APIキーが設定されていない場合はスキップ
	if os.Getenv("GEMINI_API_KEY") == "" {
		t.Skip("GEMINI_API_KEY is not set")
	}

	service, err := NewGeminiService()
	assert.NoError(t, err)

	t.Run("期限切れ間近の食材がある場合", func(t *testing.T) {
		foodItems := []model.FoodItem{
			{
				ID:         1,
				Title:      "トマト",
				Quantity:   2,
				ExpiryDate: time.Now().Add(24 * time.Hour * 3), // 3日後
			},
			{
				ID:         2,
				Title:      "なす",
				Quantity:   1,
				ExpiryDate: time.Now().Add(24 * time.Hour * 5), // 5日後
			},
		}

		recipe, err := service.GenerateRecipe(foodItems)
		assert.NoError(t, err)
		assert.NotEmpty(t, recipe)
		// レシピに食材名が含まれていることを確認
		assert.Contains(t, recipe, "トマト")
		assert.Contains(t, recipe, "なす")
	})

	t.Run("食材が空の場合", func(t *testing.T) {
		_, err := service.GenerateRecipe([]model.FoodItem{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "食材が指定されていません")
	})

	t.Run("期限切れ間近の食材がない場合", func(t *testing.T) {
		foodItems := []model.FoodItem{
			{
				ID:         1,
				Title:      "りんご",
				Quantity:   3,
				ExpiryDate: time.Now().Add(24 * time.Hour * 30), // 30日後
			},
		}

		_, err := service.GenerateRecipe(foodItems)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "期限切れ間近の食材がありません")
	})
}

func TestNewGeminiService(t *testing.T) {
	t.Run("APIキーが設定されている場合", func(t *testing.T) {
		if os.Getenv("GEMINI_API_KEY") == "" {
			t.Skip("GEMINI_API_KEY is not set")
		}

		service, err := NewGeminiService()
		assert.NoError(t, err)
		assert.NotNil(t, service)
	})

	t.Run("APIキーが設定されていない場合", func(t *testing.T) {
		// 一時的にAPIキーを削除
		originalKey := os.Getenv("GEMINI_API_KEY")
		os.Unsetenv("GEMINI_API_KEY")
		defer os.Setenv("GEMINI_API_KEY", originalKey)

		service, err := NewGeminiService()
		assert.Error(t, err)
		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "GEMINI_API_KEY")
	})
}
