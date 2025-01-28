package controllers

import (
	"net/http"
	"strconv"

	"smart-pantry/internal/services"

	"github.com/gin-gonic/gin"
)

type FoodController struct {
	foodService services.FoodService
}

func NewFoodController(foodService services.FoodService) *FoodController {
	return &FoodController{foodService: foodService}
}

// CreateFood 新しい食材を作成
func (c *FoodController) CreateFood(ctx *gin.Context) {
	var input services.FoodInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food, err := c.foodService.CreateFood(input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, food)
}

// GetFoods 食材一覧を取得
func (c *FoodController) GetFoods(ctx *gin.Context) {
	foods, err := c.foodService.GetAllFoods()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, foods)
}

// UpdateFood 食材を更新
func (c *FoodController) UpdateFood(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var input services.FoodInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	food, err := c.foodService.UpdateFood(id, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, food)
}

// DeleteFood 食材を削除
func (c *FoodController) DeleteFood(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = c.foodService.DeleteFood(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
