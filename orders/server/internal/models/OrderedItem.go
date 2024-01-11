package models

import "gorm.io/gorm"

type OrderedItem struct {
	gorm.Model
	OrderedQuantity int32
	OrderedItemId   int64
	Sku             string
	Price           float64
	OrderID         uint `gorm:"column:order_id"` // Foreign key to the Order model
}
