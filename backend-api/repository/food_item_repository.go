package repository

import (
	"backend-api/model"

	"gorm.io/gorm"
)

type IFoodItemRepository interface {
	GetAllFoodItems(foodItems *[]model.FoodItem) error
	GetFoodItemById(foodItem *model.FoodItem, id uint) error
	CreateFoodItem(foodItem *model.FoodItem) error
	UpdateFoodItem(foodItem *model.FoodItem) error
	DeleteFoodItem(id uint) error
}

type foodItemRepository struct {
	db *gorm.DB
}

func NewFoodItemRepository(db *gorm.DB) IFoodItemRepository {
	return &foodItemRepository{db}
}

func (fr *foodItemRepository) GetAllFoodItems(foodItems *[]model.FoodItem) error {
	if err := fr.db.Find(foodItems).Error; err != nil {
		return err
	}
	return nil
}

func (fr *foodItemRepository) GetFoodItemById(foodItem *model.FoodItem, id uint) error {
	if err := fr.db.First(foodItem, id).Error; err != nil {
		return err
	}
	return nil
}

func (fr *foodItemRepository) CreateFoodItem(foodItem *model.FoodItem) error {
	if err := fr.db.Create(foodItem).Error; err != nil {
		return err
	}
	return nil
}

func (fr *foodItemRepository) UpdateFoodItem(foodItem *model.FoodItem) error {
	if err := fr.db.Save(foodItem).Error; err != nil {
		return err
	}
	return nil
}

func (fr *foodItemRepository) DeleteFoodItem(id uint) error {
	if err := fr.db.Delete(&model.FoodItem{}, id).Error; err != nil {
		return err
	}
	return nil
}
