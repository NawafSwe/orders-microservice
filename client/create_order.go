package main

import (
	"context"
	"github.com/nawafswe/orders-service/proto"
	"log"
	"time"
)

func createOrder(c proto.OrderServiceClient) {
	items := []*proto.OrderedItem{
		{
			OrderedItemId:   int64(24),
			OrderedQuantity: 1,
			Price:           10.00,
			Sku:             "XYZ123456",
		},
	}
	req := &proto.Order{
		CustomerId: 1,
		GrandTotal: 10.00,
		Items:      items,
		Status:     "New",
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := c.Create(ctx, req)

	if err != nil {
		log.Fatalf("Error occurred when calling create order service, err: %v\n", err)
	}

	log.Printf("created order -> %v\n", res)

}