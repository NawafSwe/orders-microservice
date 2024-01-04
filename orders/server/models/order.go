package models

type Order struct {
	OrderId    string `gorm:"primarykey"`
	CustomerId string
	Status     string
	GrandTotal float64
	// CreatedAt time.Time `gorm:"autoCreateTime:false"`
	// UpdatedAt time.Time `gorm:"autoUpdateTime:false"`
}

// message Order {

//     int64 order_id = 1;
//     int64 customer_id = 2;
//     string status = 3;
//     double grand_total = 4;
//     repeated OrderedItem items = 5;

// }

// message OrderStatus {
//     int64 order_id = 1;
//     string status = 2;

// }
