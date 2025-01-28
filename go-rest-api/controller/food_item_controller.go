package controller

import (
	"go-rest-api/model"
	"go-rest-api/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type IFoodItemController interface {
	GetAllFoodItems(c echo.Context) error
	GetFoodItemById(c echo.Context) error
	CreateFoodItem(c echo.Context) error
	UpdateFoodItem(c echo.Context) error
	DeleteFoodItem(c echo.Context) error
}

type foodItemController struct {
	fu usecase.IFoodItemUsecase
}

func NewFoodItemController(fu usecase.IFoodItemUsecase) IFoodItemController {
	return &foodItemController{fu}
}

func (fc *foodItemController) GetAllFoodItems(c echo.Context) error {
	foodItems, err := fc.fu.GetAllFoodItems()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, foodItems)
}

func (fc *foodItemController) GetFoodItemById(c echo.Context) error {
	id := c.Param("id")
	foodItemId, _ := strconv.Atoi(id)
	foodItem, err := fc.fu.GetFoodItemById(uint(foodItemId))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, foodItem)
}

func (fc *foodItemController) CreateFoodItem(c echo.Context) error {
	foodItem := model.FoodItem{}
	if err := c.Bind(&foodItem); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	createdFoodItem, err := fc.fu.CreateFoodItem(foodItem)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, createdFoodItem)
}

func (fc *foodItemController) UpdateFoodItem(c echo.Context) error {
	foodItem := model.FoodItem{}
	if err := c.Bind(&foodItem); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := fc.fu.UpdateFoodItem(foodItem); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, foodItem)
}

func (fc *foodItemController) DeleteFoodItem(c echo.Context) error {
	id := c.Param("id")
	foodItemId, _ := strconv.Atoi(id)
	if err := fc.fu.DeleteFoodItem(uint(foodItemId)); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}
