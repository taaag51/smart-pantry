package services

import (
	"context"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type RecipeService struct {
	geminiClient *genai.Client
	geminiModel  *genai.GenerativeModel
}

func NewRecipeService() (*RecipeService, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-pro")

	return &RecipeService{
		geminiClient: client,
		geminiModel:  model,
	}, nil
}

func (rs *RecipeService) SuggestRecipe(ingredients []string) (string, error) {
	prompt := "冷蔵庫にある食材" + strings.Join(ingredients, ",") + "を使って作れる簡単な料理を提案してください。"

	resp, err := rs.geminiModel.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		recipe := resp.Candidates[0].Content.Parts[0].(genai.Text).String()
		return recipe, nil
	}

	return "レシピの提案に失敗しました。", nil
}

func (rs *RecipeService) Close() error {
	return rs.geminiClient.Close()
}
