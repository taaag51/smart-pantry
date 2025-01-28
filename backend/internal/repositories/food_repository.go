package repositories

import (
	"smart-pantry/internal/models"

	"gorm.io/gorm"
)

type FoodRepository interface {
	Create(food *models.FoodItem) error
	FindAll() ([]models.FoodItem, error)
	FindByID(id int) (*models.FoodItem, error)
	Update(food *models.FoodItem) error
	Delete(id int) error
}

type foodRepository struct {
	db *gorm.DB
}

func NewFoodRepository(db *gorm.DB) FoodRepository {
	return &foodRepository{db: db}
}

func (r *foodRepository) Create(food *models.FoodItem) error {
	return r.db.Create(food).Error
}

func (r *foodRepository) FindAll() ([]models.FoodItem, error) {
	var foods []models.FoodItem
	if err := r.db.Find(&foods).Error; err != nil {
		return nil, err
	}
	return foods, nil
}

func (r *foodRepository) FindByID(id int) (*models.FoodItem, error) {
	var food models.FoodItem
	if err := r.db.First(&food, id).Error; err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *foodRepository) Update(food *models.FoodItem) error {
	return r.db.Save(food).Error
}

func (r *foodRepository) Delete(id int) error {
	return r.db.Delete(&models.FoodItem{}, id).Error
}
