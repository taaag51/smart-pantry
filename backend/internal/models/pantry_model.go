package models

import "gorm.io/gorm"

// PantryItem represents the pantry_items table in the database.
type PantryItem struct {
	gorm.Model
	Name       string `gorm:"not null"`
	Quantity   int    `gorm:"not null"`
	ExpiryDate string `gorm:"type:date"` // Consider using time.Time and proper formatting
}
