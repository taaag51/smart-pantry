package services

import (
	"smart-pantry/internal/models"
	"smart-pantry/internal/repositories"
)

type FoodService interface {
	CreateFood(input FoodInput) (*models.FoodItem, error)
	GetAllFoods() ([]models.FoodItem, error)
	UpdateFood(id int, input FoodInput) (*models.FoodItem, error)
	DeleteFood(id int) error
}

type foodService struct {
	repo repositories.FoodRepository
}

type FoodInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
	Unit        string `json:"unit"`
	ExpiryDate  string `json:"expiry_date"`
}

func NewFoodService(repo repositories.FoodRepository) FoodService {
	return &foodService{repo: repo}
}

func (s *foodService) CreateFood(input FoodInput) (*models.FoodItem, error) {
	food := models.FoodItem{
		Name:        input.Name,
		Description: input.Description,
		Quantity:    input.Quantity,
		Unit:        input.Unit,
		ExpiryDate:  input.ExpiryDate,
	}

	if err := s.repo.Create(&food); err != nil {
		return nil, err
	}

	return &food, nil
}

func (s *foodService) GetAllFoods() ([]models.FoodItem, error) {
	return s.repo.FindAll()
}

func (s *foodService) UpdateFood(id int, input FoodInput) (*models.FoodItem, error) {
	food, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	food.Name = input.Name
	food.Description = input.Description
	food.Quantity = input.Quantity
	food.Unit = input.Unit
	food.ExpiryDate = input.ExpiryDate

	if err := s.repo.Update(food); err != nil {
		return nil, err
	}

	return food, nil
}

func (s *foodService) DeleteFood(id int) error {
	return s.repo.Delete(id)
}
