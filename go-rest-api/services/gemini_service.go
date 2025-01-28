package services

import (
	"context"
	"fmt"
	"go-rest-api/model"
	"os"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type IGeminiService interface {
	GenerateRecipe(foodItems []model.FoodItem) (string, error)
}

type geminiService struct {
	client *genai.Client
	model  *genai.GenerativeModel
}

func NewGeminiService() (IGeminiService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, fmt.Errorf("Gemini APIクライアントの作成に失敗しました: %v", err)
	}

	model := client.GenerativeModel("gemini-pro")
	return &geminiService{
		client: client,
		model:  model,
	}, nil
}

func (s *geminiService) GenerateRecipe(foodItems []model.FoodItem) (string, error) {
	if len(foodItems) == 0 {
		return "", fmt.Errorf("食材が指定されていません")
	}

	// 期限切れ間近の食材を抽出
	var expiringItems []model.FoodItem
	for _, item := range foodItems {
		daysUntilExpiry := time.Until(item.ExpiryDate).Hours() / 24
		if daysUntilExpiry > 0 && daysUntilExpiry <= 7 {
			expiringItems = append(expiringItems, item)
		}
	}

	if len(expiringItems) == 0 {
		return "", fmt.Errorf("期限切れ間近の食材がありません")
	}

	// プロンプトの構築
	var promptBuilder strings.Builder
	promptBuilder.WriteString("以下の食材を使用した、栄養バランスの良いレシピを提案してください：\n\n")
	promptBuilder.WriteString("【食材リスト】\n")
	for _, item := range expiringItems {
		promptBuilder.WriteString(fmt.Sprintf("- %s（%d個）: 賞味期限 %s\n",
			item.Title,
			item.Quantity,
			item.ExpiryDate.Format("2006/01/02")))
	}
	promptBuilder.WriteString("\n【条件】\n")
	promptBuilder.WriteString("1. 上記の食材を優先的に使用すること\n")
	promptBuilder.WriteString("2. 栄養バランスを考慮すること\n")
	promptBuilder.WriteString("3. 調理手順は簡潔に記載すること\n")
	promptBuilder.WriteString("4. 必要な追加食材があれば提案すること\n")
	promptBuilder.WriteString("\n【出力形式】\n")
	promptBuilder.WriteString("1. レシピ名\n")
	promptBuilder.WriteString("2. 材料（2人分）\n")
	promptBuilder.WriteString("3. 調理手順\n")
	promptBuilder.WriteString("4. 栄養バランスの説明\n")

	// Gemini APIにリクエスト
	ctx := context.Background()
	resp, err := s.model.GenerateContent(ctx, genai.Text(promptBuilder.String()))
	if err != nil {
		return "", fmt.Errorf("レシピの生成に失敗しました: %v", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("レシピを生成できませんでした")
	}

	// レスポンスの取得
	recipe, ok := resp.Candidates[0].Content.Parts[0].(genai.Text)
	if !ok {
		return "", fmt.Errorf("レスポンスの形式が不正です")
	}

	return string(recipe), nil
}
