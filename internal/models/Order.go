package models

import "gorm.io/gorm"

type OrderStatus int

// iota starts from zero and it increment by one for each
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
	CustomerId int64
	Status     string
	GrandTotal float64
	Items      []OrderedItem `gorm:"foreignKey:order_id"` // one to many
}
