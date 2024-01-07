package models

import "gorm.io/gorm"

type OrderedItem struct {
	gorm.Model
	OrderedQuantity int32
	Sku             string
	Price           float64
	OrderID         uint `gorm:"column:id"` // Foreign key to the Order model
}
