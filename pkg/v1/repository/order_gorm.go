package repo

import (
	"github.com/nawafswe/orders-service/internal/models"
	interfaces "github.com/nawafswe/orders-service/pkg/v1"
	"gorm.io/gorm"
)

type OrderRepoImpl struct {
	db *gorm.DB
}

func New(d *gorm.DB) interfaces.OrderRepo {
	return OrderRepoImpl{db: d}
}

func (r OrderRepoImpl) Create(order models.Order) (models.Order, error) {
	return models.Order{}, nil
}
func (r OrderRepoImpl) UpdateOrderStatus(id int64, status string) (models.Order, error) {
	return models.Order{}, nil
}
