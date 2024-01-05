package domain

import (
	"context"
	"errors"
	"log"

	"github.com/nawafswe/orders-service/orders/server/models"
	"gorm.io/gorm"
)

func ChangeOrderStatus(ctx context.Context, o models.Order, session *gorm.DB) error {
	log.Printf("ChangeOrderStatus was invoked with o: %v\n", o)
	var order models.Order
	tx := session.WithContext(ctx).Model(&order).Where("order_id=?", o.OrderId).Updates(models.Order{Status: o.Status})

	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("order not found")
	}

	return nil

}
