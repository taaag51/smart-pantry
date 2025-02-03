package usecase

import (
	"github.com/taaag51/smart-pantry/backend-api/model"
	"github.com/taaag51/smart-pantry/backend-api/repository"
)

type IFoodItemUsecase interface {
	GetAllFoodItems() ([]model.FoodItem, error)
	GetFoodItemById(id uint) (model.FoodItem, error)
	CreateFoodItem(foodItem model.FoodItem) (model.FoodItem, error)
	UpdateFoodItem(foodItem model.FoodItem) error
	DeleteFoodItem(id uint) error
}

type foodItemUsecase struct {
	fr repository.IFoodItemRepository
}

func NewFoodItemUsecase(fr repository.IFoodItemRepository) IFoodItemUsecase {
	return &foodItemUsecase{fr}
}

func (fu *foodItemUsecase) GetAllFoodItems() ([]model.FoodItem, error) {
	foodItems := []model.FoodItem{}
	if err := fu.fr.GetAllFoodItems(&foodItems); err != nil {
		return nil, err
	}
	return foodItems, nil
}

func (fu *foodItemUsecase) GetFoodItemById(id uint) (model.FoodItem, error) {
	foodItem := model.FoodItem{}
	if err := fu.fr.GetFoodItemById(&foodItem, id); err != nil {
		return model.FoodItem{}, err
	}
	return foodItem, nil
}

func (fu *foodItemUsecase) CreateFoodItem(foodItem model.FoodItem) (model.FoodItem, error) {
	if err := fu.fr.CreateFoodItem(&foodItem); err != nil {
		return model.FoodItem{}, err
	}
	return foodItem, nil
}

func (fu *foodItemUsecase) UpdateFoodItem(foodItem model.FoodItem) error {
	if err := fu.fr.UpdateFoodItem(&foodItem); err != nil {
		return err
	}
	return nil
}

func (fu *foodItemUsecase) DeleteFoodItem(id uint) error {
	if err := fu.fr.DeleteFoodItem(id); err != nil {
		return err
	}
	return nil
}
