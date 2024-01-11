package v1

import (
	"github.com/nawafswe/orders-service/internal/models"
)

type OrderRepo interface {
	Create(order models.Order) (models.Order, error)
	UpdateOrderStatus(id int64, status string) (models.Order, error)
}

type OrderUseCase interface {
	PlaceOrder(order models.Order) (models.Order, error)
	ApproveOrder(orderId int64) error
	RejectOrder(orderId int64) error
}
