package repo

import (
	"context"
	"errors"
	"fmt"
	interfaces "github.com/nawafswe/orders-service/internal/app/orders"
	"github.com/nawafswe/orders-service/internal/models"
	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	db *gorm.DB
}

func NewOrderRepo(d *gorm.DB) interfaces.OrderRepo {
	return OrderRepoImpl{db: d}
}

func (r OrderRepoImpl) Create(ctx context.Context, order models.Order) (models.Order, error) {
	tx := r.db.WithContext(ctx).Create(&order)

	if tx.Error != nil {
		return models.Order{}, fmt.Errorf("error occurred while creating a new order, err: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return models.Order{}, errors.New("failed to create order for unknown error")
	}
	return order, nil
}
func (r OrderRepoImpl) UpdateOrderStatus(ctx context.Context, id int64, status string) (models.Order, error) {

	tx := r.db.WithContext(ctx).Where("id = ? ", id).Updates(&models.Order{Status: status})
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return models.Order{}, fmt.Errorf("updating status failed due to invalid order id was given %d", id)
		}
		return models.Order{}, fmt.Errorf("UpdateOrderStatus: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return models.Order{}, fmt.Errorf("order with id %v not found", id)
	}
	var o models.Order
	tx.Scan(&o)

	return o, nil
}
