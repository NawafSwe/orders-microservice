package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/nawafswe/orders-service/proto"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

func createOrder(c proto.OrderServiceClient) {
	items := []*proto.OrderedItem{
		{
			OrderedItemId:   int64(12),
			OrderedQuantity: 1,
			Price:           25,
			Name:            "V60",
		},
	}
	req := &proto.Order{
		CustomerId:   1,
		RestaurantId: 33,
		GrandTotal:   30,
		Items:        items,
		Status:       "New",
	}

	md := metadata.Pairs("correlation-id", uuid.New().String())
	ctx, cancel := context.WithTimeout(metadata.NewOutgoingContext(context.Background(), md), time.Second*10)
	defer cancel()
	res, err := c.Create(ctx, req)

	if err != nil {
		log.Fatalf("Error occurred when calling create order service, err: %v\n", err)
	}

	log.Printf("created order -> %v\n", res)

}
