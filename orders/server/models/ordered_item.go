package models

import "gorm.io/gorm"

type OrderedItem struct {
	gorm.Model
	ItemId          int64 `gorm:"uniqueIndex;autoIncrement"`
	OrderedQuantity int32
	Sku             string
	Price           float64
	OrderID         uint `gorm:"column:order_id"` // Foreign key to the Order model
}
