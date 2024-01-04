package models

type Order struct {
	OrderId    string `gorm:"primarykey"`
	CustomerId string
	Status     string
	GrandTotal float64
	// CreatedAt time.Time `gorm:"autoCreateTime:false"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}
