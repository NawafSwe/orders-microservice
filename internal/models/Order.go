package models

import (
	"gorm.io/gorm"
)

type OrderStatus int

// iota starts from zero, and it increments by one for each
// you can manually set it for other fields
// once this block is created iota is set back to zero
const (
	New OrderStatus = iota
	Approved
	Rejected
	Cancelled
	Delivered
)

type Order struct {
	gorm.Model
	CustomerId   int64
	RestaurantId int64
	Status       string
	GrandTotal   float64
	Items        []OrderedItem `gorm:"foreignKey:order_id"` // one to many
}

type InvalidStatusChangeErr struct {
	Message string
}

func (i InvalidStatusChangeErr) Error() string {
	return i.Message
}
