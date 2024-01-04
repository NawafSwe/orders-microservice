package main

import (
	"context"
	"log"
	"time"

	pb "github.com/nawafswe/orders-service/orders/proto"
)

func createOrder(c pb.OrderServiceClient) {
	items := []*pb.OrderedItem{
		{
			OrderedQuantity: 1,
			Price:           45.00,
			Sku:             "drink-coffee-321-s",
		},
	}
	req := &pb.Order{
		CustomerId: 1,
		GrandTotal: 45.00,
		Items:      items,
		Status:     "New",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	res, err := c.Create(ctx, req)

	if err != nil {
		log.Fatalf("Error occurred when calling create order service, err: %v\n", err)
	}

	log.Printf("created order -> %v\n", res)

}
