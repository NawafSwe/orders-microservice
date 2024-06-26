package main

import (
	"context"
	"github.com/nawafswe/orders-service/proto"
	"log"
	"time"
)

func changeOrderStatus(c proto.OrderServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := c.ChangeOrderStatus(ctx, &proto.OrderStatus{
		OrderId: 1,
		Status:  "Under-Preparation",
	})

	if err != nil {
		log.Fatalf("failed to update order status, err: %v\n", err)
	}

	log.Printf("successfully updated order status")
}
