package models

import "gorm.io/gorm"

type OrderedItem struct {
	gorm.Model
	ItemId          int64 `gorm:"primarykey"`
	OrderedQuantity int32
	Sku             string
	price           float64
	OrderID         uint `gorm:"column:order_id"` // Foreign key to the Order model
}
