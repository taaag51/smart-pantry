package services

import (
	"testing"
	"time"

	"backend-api/model"
	"backend-api/testutil"

	"github.com/stretchr/testify/assert"
)

func TestGeminiService_GenerateRecipe(t *testing.T) {
	// テストデータベースのセットアップ
	db := testutil.NewTestDB(t)
	defer db.Close()

	// テストユーザーの作成
	userID, _ := db.CreateTestUser(t)

	// テストケース
	tests := []struct {
		name        string
		ingredients []model.FoodItem
		wantErr     bool
		errMessage  string
	}{
		{
			name:        "空の食材リスト",
			ingredients: []model.FoodItem{},
			wantErr:     true,
			errMessage:  "食材が指定されていません",
		},
		{
			name: "有効な食材リスト",
			ingredients: []model.FoodItem{
				{
					ID:         1,
					UserId:     userID,
					Title:      "玉ねぎ",
					Quantity:   2,
					ExpiryDate: time.Now().AddDate(0, 0, 7),
				},
				{
					ID:         2,
					UserId:     userID,
					Title:      "じゃがいも",
					Quantity:   3,
					ExpiryDate: time.Now().AddDate(0, 0, 7),
				},
			},
			wantErr: false,
		},
		{
			name: "期限切れ食材を含む",
			ingredients: []model.FoodItem{
				{
					ID:         3,
					UserId:     userID,
					Title:      "期限切れ食材",
					Quantity:   1,
					ExpiryDate: time.Now().AddDate(0, 0, -1), // 昨日で期限切れ
				},
			},
			wantErr: false, // 期限切れ食材も受け付けるように変更
		},
	}

	// GeminiServiceのインスタンス作成
	service, err := NewGeminiService()
	if err != nil {
		t.Fatalf("GeminiServiceの作成に失敗: %v", err)
	}

	// テストケースの実行
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recipe, err := service.GenerateRecipe(tt.ingredients)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, recipe)
				// レシピの内容に関する基本的な検証
				assert.Contains(t, recipe, "作り方")
				assert.Contains(t, recipe, "材料")
			}
		})
	}
}

func TestGeminiService_GenerateRecipe_LargeQuantity(t *testing.T) {
	// テストデータベースのセットアップ
	db := testutil.NewTestDB(t)
	defer db.Close()

	// テストユーザーの作成
	userID, _ := db.CreateTestUser(t)

	// 大量の食材データを生成
	var ingredients []model.FoodItem
	for i := 1; i <= 50; i++ {
		ingredients = append(ingredients, model.FoodItem{
			ID:         uint(i),
			UserId:     userID,
			Title:      "テスト食材",
			Quantity:   i,
			ExpiryDate: time.Now().AddDate(0, 0, 7),
		})
	}

	// GeminiServiceのインスタンス作成
	service, err := NewGeminiService()
	if err != nil {
		t.Fatalf("GeminiServiceの作成に失敗: %v", err)
	}

	// テスト実行
	recipe, err := service.GenerateRecipe(ingredients)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipe)
	assert.Contains(t, recipe, "大量調理")
}

func TestGeminiService_GenerateRecipe_Performance(t *testing.T) {
	// テストデータベースのセットアップ
	db := testutil.NewTestDB(t)
	defer db.Close()

	// テストユーザーの作成
	userID, _ := db.CreateTestUser(t)

	// 標準的な食材リスト
	ingredients := []model.FoodItem{
		{
			ID:         1,
			UserId:     userID,
			Title:      "玉ねぎ",
			Quantity:   2,
			ExpiryDate: time.Now().AddDate(0, 0, 7),
		},
		{
			ID:         2,
			UserId:     userID,
			Title:      "じゃがいも",
			Quantity:   3,
			ExpiryDate: time.Now().AddDate(0, 0, 7),
		},
	}

	// GeminiServiceのインスタンス作成
	service, err := NewGeminiService()
	if err != nil {
		t.Fatalf("GeminiServiceの作成に失敗: %v", err)
	}

	// パフォーマンステスト
	start := time.Now()
	recipe, err := service.GenerateRecipe(ingredients)
	duration := time.Since(start)

	// アサーション
	assert.NoError(t, err)
	assert.NotEmpty(t, recipe)
	assert.Less(t, duration.Seconds(), 2.0, "レシピ生成は2秒以内に完了すべき")
}
