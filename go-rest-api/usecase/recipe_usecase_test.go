package usecase

import (
	"go-rest-api/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFoodItemRepository struct {
	mock.Mock
}

func (m *MockFoodItemRepository) GetAllFoodItems(foodItems *[]model.FoodItem) error {
	args := m.Called(foodItems)
	return args.Error(0)
}

type MockGeminiService struct {
	mock.Mock
}

func (m *MockGeminiService) GenerateRecipe(foodItems []model.FoodItem) (string, error) {
	args := m.Called(foodItems)
	return args.String(0), args.Error(1)
}

func TestGetRecipeSuggestions(t *testing.T) {
	mockRepo := new(MockFoodItemRepository)
	mockGemini := new(MockGeminiService)
	usecase := NewRecipeUsecase(mockRepo, mockGemini)

	t.Run("期限切れ間近の食材がある場合", func(t *testing.T) {
		// テストデータ
		foodItems := []model.FoodItem{
			{
				ID:         1,
				Title:      "トマト",
				Quantity:   2,
				ExpiryDate: time.Now().Add(24 * time.Hour * 3), // 3日後
			},
		}
		expectedRecipe := "トマトを使用したレシピ..."

		// モックの設定
		mockRepo.On("GetAllFoodItems", mock.AnythingOfType("*[]model.FoodItem")).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*[]model.FoodItem)
				*arg = foodItems
			}).
			Return(nil)

		mockGemini.On("GenerateRecipe", foodItems).
			Return(expectedRecipe, nil)

		// テスト実行
		recipe, err := usecase.GetRecipeSuggestions(1)

		// アサーション
		assert.NoError(t, err)
		assert.Equal(t, expectedRecipe, recipe)
		mockRepo.AssertExpectations(t)
		mockGemini.AssertExpectations(t)
	})

	t.Run("食材が存在しない場合", func(t *testing.T) {
		// テストデータ
		var emptyFoodItems []model.FoodItem

		// モックの設定
		mockRepo.On("GetAllFoodItems", mock.AnythingOfType("*[]model.FoodItem")).
			Run(func(args mock.Arguments) {
				arg := args.Get(0).(*[]model.FoodItem)
				*arg = emptyFoodItems
			}).
			Return(nil)

		// テスト実行
		recipe, err := usecase.GetRecipeSuggestions(1)

		// アサーション
		assert.Error(t, err)
		assert.Empty(t, recipe)
		mockRepo.AssertExpectations(t)
	})
}
