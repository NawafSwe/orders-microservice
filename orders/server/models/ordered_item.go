package models

type OrderedItem struct {
	ItemId          string `gorm:"primarykey"`
	OrderedQuantity int32
	Sku             string
	price           float64
}
