package models

type PantryItem struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Quantity int    `gorm:"not null"`
}
