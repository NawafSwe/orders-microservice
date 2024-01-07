package domain

import (
	"context"
	"errors"
	"log"

	"github.com/nawafswe/orders-service/orders/server/internal/models"
	"gorm.io/gorm"
)

func ChangeOrderStatus(ctx context.Context, orderId int64, o models.Order, session *gorm.DB) error {
	log.Printf("ChangeOrderStatus was invoked with o: %v\n", o)
	var order models.Order
	tx := session.WithContext(ctx).Model(&order).Where("id=?", orderId).Updates(models.Order{Status: o.Status})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil

}
