package usecase

import (
	"github.com/nawafswe/orders-service/internal/models"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
)

type OrderUseCaseImpl struct {
	repo interfaces.OrderRepo
}

func New(repo interfaces.OrderRepo) interfaces.OrderUseCase {
	return OrderUseCaseImpl{repo: repo}
}

func (ou OrderUseCaseImpl) PlaceOrder(order models.Order) (models.Order, error) {
	return models.Order{}, nil
}

func (ou OrderUseCaseImpl) ApproveOrder(orderId int64) error {
	return nil
}
func (ou OrderUseCaseImpl) RejectOrder(orderId int64) error {
	return nil
}
