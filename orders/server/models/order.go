package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderId    int64 `gorm:"uniqueIndex;autoIncrement"`
	CustomerId int64
	Status     string
	GrandTotal float64
	Items      []OrderedItem `gorm:"foreignKey:ItemId"` // one to many
	// CreatedAt time.Time `gorm:"autoCreateTime:false"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
