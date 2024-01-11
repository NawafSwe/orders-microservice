package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	CustomerId int64
	Status     string
	GrandTotal float64
	Items      []OrderedItem `gorm:"foreignKey:order_id"` // one to many
}
