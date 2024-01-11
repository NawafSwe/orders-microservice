package v1

import (
	"context"
	"github.com/nawafswe/orders-service/internal/models"
)

type OrderRepo interface {
	Create(ctx context.Context, order models.Order) (models.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status string) (models.Order, error)
}

type OrderUseCase interface {
	PlaceOrder(ctx context.Context, order models.Order) (models.Order, error)
	ApproveOrder(ctx context.Context, orderId int64) error
	RejectOrder(ctx context.Context, orderId int64) error
	UpdateOrderStatus(ctx context.Context, orderId int64, status string) (models.Order, error)
	HandleOrderApproval(ctx context.Context)
}
