package services

import (
	"smart-pantry/backend/internal/models"

	"gorm.io/gorm"
)

type PantryService struct {
	db *gorm.DB
}

func NewPantryService(db *gorm.DB) *PantryService {
	return &PantryService{db: db}
}

type PantryItemParams struct {
	Name       string  `json:"name"`
	Quantity   int     `json:"quantity"`
	ExpiryDate *string `json:"expiry_date"`
}

func (ps *PantryService) GetPantryItems() ([]models.PantryItem, error) {
	var items []models.PantryItem
	result := ps.db.Find(&items)
	return items, result.Error
}

func (ps *PantryService) CreatePantryItem(params PantryItemParams) (*models.PantryItem, error) {
	item := models.PantryItem{
		Name:       params.Name,
		Quantity:   params.Quantity,
		ExpiryDate: *params.ExpiryDate,
	}
	result := ps.db.Create(&item)
	return &item, result.Error
}

func (ps *PantryService) UpdatePantryItem(id int, params PantryItemParams) (*models.PantryItem, error) {
	var item models.PantryItem
	if err := ps.db.First(&item, id).Error; err != nil {
		return nil, err
	}

	item.Name = params.Name
	item.Quantity = params.Quantity
	item.ExpiryDate = *params.ExpiryDate

	result := ps.db.Save(&item)
	return &item, result.Error
}

func (ps *PantryService) DeletePantryItem(id int) error {
	var item models.PantryItem
	result := ps.db.Delete(&item, id)
	return result.Error
}
