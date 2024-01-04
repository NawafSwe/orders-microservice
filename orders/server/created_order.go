package main

import (
	"context"
	"fmt"

	pb "github.com/nawafswe/orders-service/orders/proto"
)

func (s *Server) Create(ctx context.Context, in *pb.Order) (*pb.Order, error) {

	fmt.Printf("OrderService was invoked with Create method, with ctx: %v, in:%v\n", ctx, in)
	items := []*pb.OrderedItem{
		{
			ItemId:          1,
			OrderedQuantity: 1,
			Price:           45.00,
			Sku:             "drink-coffee-321-s",
		},
	}
	return &pb.Order{
		OrderId:    1,
		CustomerId: 1,
		GrandTotal: 45.00,
		Items:      items,
		Status:     "New",
	}, nil
}
