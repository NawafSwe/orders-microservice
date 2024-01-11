package usecase

import (
	"context"
	"github.com/nawafswe/orders-service/internal/models"
	"github.com/nawafswe/orders-service/pkg/messaging"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
	"github.com/nawafswe/orders-service/pkg/v1/handler/grpc"
)

type OrderUseCaseImpl struct {
	repo interfaces.OrderRepo
	// define an interface for messaging once you segregate business logic for publishing events and handling events from there
	pubSubClient messaging.PUBSUB
}

func New(repo interfaces.OrderRepo, ps messaging.PUBSUB) interfaces.OrderUseCase {
	return OrderUseCaseImpl{repo: repo, pubSubClient: ps}
}

func (ou OrderUseCaseImpl) PlaceOrder(ctx context.Context, order models.Order) (models.Order, error) {
	o, err := ou.repo.Create(ctx, order)
	if err != nil {
		return order, err
	}
	ou.pubSubClient.PublishOrderCreatedEvent(ctx, grpc.FromDomain(order))
	return o, nil
}

func (ou OrderUseCaseImpl) ApproveOrder(ctx context.Context, orderId int64) error {
	return nil
}
func (ou OrderUseCaseImpl) RejectOrder(ctx context.Context, orderId int64) error {
	return nil
}
func (ou OrderUseCaseImpl) UpdateOrderStatus(ctx context.Context, orderId int64, status string) (models.Order, error) {
	return ou.repo.UpdateOrderStatus(ctx, orderId, status)
}
