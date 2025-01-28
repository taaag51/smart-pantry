package models

type FoodItem struct {
	ID          int    `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description" gorm:"type:text"`
	Quantity    int    `json:"quantity" gorm:"not null"`
	Unit        string `json:"unit" gorm:"size:50"`
	ExpiryDate  string `json:"expiry_date" gorm:"type:date"`
}
