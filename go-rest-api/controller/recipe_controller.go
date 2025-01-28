package controller

import (
	"go-rest-api/usecase"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type IRecipeController interface {
	GetRecipeSuggestions(c echo.Context) error
}

type recipeController struct {
	ru usecase.IRecipeUsecase
}

func NewRecipeController(ru usecase.IRecipeUsecase) IRecipeController {
	return &recipeController{ru}
}

func (rc *recipeController) GetRecipeSuggestions(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userId := uint(claims["user_id"].(float64))

	// レシピ提案を取得
	suggestions, err := rc.ru.GetRecipeSuggestions(userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "レシピの提案に失敗しました",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"recipe": suggestions,
	})
}
