package model

import "time"

// FoodItem represents a food item with its details.
type FoodItem struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Title      string    `json:"title" gorm:"not null"`       // Reusing the Title field from Task
	Quantity   int       `json:"quantity" gorm:"not null"`    // New field for quantity
	ExpiryDate time.Time `json:"expiry_date" gorm:"not null"` // New field for expiry date
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId     uint      `json:"user_id" gorm:"not null"`
}

// FoodItemResponse is the response structure for food items
type FoodItemResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	Quantity   int       `json:"quantity"`
	ExpiryDate time.Time `json:"expiry_date"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
