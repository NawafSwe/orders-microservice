package main

import (
	"context"
	"log"
	"time"

	pb "github.com/nawafswe/orders-service/orders/proto"
)

func changeOrderStatus(c pb.OrderServiceClient) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	_, err := c.ChangeOrderStatus(ctx, &pb.OrderStatus{
		OrderId: 3,
		Status:  "Confirmed",
	})

	if err != nil {
		log.Fatalf("failed to update order status, err: %v\n", err)
	}

	log.Printf("successfully updated order status")
}
